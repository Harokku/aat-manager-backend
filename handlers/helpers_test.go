package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"net/url"
	"testing"
)

func TestGetRedirectPath(t *testing.T) {
	tests := []struct {
		name           string
		page           string
		ctxPath        string
		expectRedirect string
	}{
		{
			name:           "Empty Page and Path",
			page:           "",
			ctxPath:        "",
			expectRedirect: "?redirect=",
		},
		{
			name:           "Non-Empty Page, Empty Path",
			page:           "/home",
			ctxPath:        "",
			expectRedirect: "/home?redirect=",
		},
		{
			name:           "Non-Empty Page and Path",
			page:           "/home",
			ctxPath:        "/login",
			expectRedirect: "/home?redirect=/login",
		},
		{
			name:           "Empty Page, Non-Empty Path",
			page:           "",
			ctxPath:        "/login",
			expectRedirect: "?redirect=/login",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(ctx)
			ctx.Path(tc.ctxPath)
			redirect := getRedirectPath(tc.page, ctx)
			if redirect != tc.expectRedirect {
				t.Errorf("Test name: %s, Expected: %s, Got: %s", tc.name, tc.expectRedirect, redirect)
			}
		})
	}
}

func TestGetChainableRedirectPath(t *testing.T) {
	testCases := []struct {
		name         string
		page         string
		redirectPage string
		expected     string
	}{
		{
			name:         "RedirectToHomePage",
			page:         "/home",
			redirectPage: "/login",
			expected:     "/home?redirect=/login",
		},
		{
			name:         "RedirectToProfilePage",
			page:         "/profile",
			redirectPage: "/home",
			expected:     "/profile?redirect=/home",
		},
		{
			name:         "NoRedirectPage",
			page:         "/profile",
			redirectPage: "",
			expected:     "/profile?redirect=",
		},
		{
			name:         "RedirectToURLWithQueryParameters",
			page:         "/search",
			redirectPage: "/products?productID=1234",
			expected:     "/search?redirect=/products?productID=1234",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app := fiber.New()
			fasthttpCtx := &fasthttp.RequestCtx{}
			fasthttpCtx.Request.SetRequestURI(fmt.Sprintf("/dummy?redirect=%s", url.QueryEscape(tc.redirectPage)))
			ctx := app.AcquireCtx(fasthttpCtx)
			defer app.ReleaseCtx(ctx)
			got := getChainableRedirectPath(tc.page, ctx)

			if got != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, got)
			}
		})
	}
}
