package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestPerformTrialCall_Success simulates a successful trial claim.
func TestPerformTrialCall_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/billing/ide/trial" {
			t.Errorf("expected /billing/ide/trial, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"msg":"OK","data":{"trial":true,"credits":250}}`))
	}))
	defer srv.Close()

	restore := setBillingBaseGlobal(srv.URL)
	defer restore()

	sa := &storedAuth{Auth: storedTokens{Domain: "www.workbuddy.ai"}}
	res, err := performTrialCall(sa)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["success"] != true {
		t.Errorf("expected success=true, got %v", res["success"])
	}
}

// TestPerformTrialCall_AlreadyClaimed simulates the code=14051 "has applied trial" response.
func TestPerformTrialCall_AlreadyClaimed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":14051,"msg":"has applied trial","data":null}`))
	}))
	defer srv.Close()

	restore := setBillingBaseGlobal(srv.URL)
	defer restore()

	sa := &storedAuth{Auth: storedTokens{Domain: "www.workbuddy.ai"}}
	res, err := performTrialCall(sa)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["success"] != false {
		t.Errorf("expected success=false, got %v", res["success"])
	}
	if res["already_claimed"] != true {
		t.Errorf("expected already_claimed=true, got %v", res["already_claimed"])
	}
	if msg, ok := res["message"].(string); !ok || !strings.Contains(msg, "已领取") {
		t.Errorf("expected message to contain 已领取, got %v", res["message"])
	}
}

// TestHasTrialPack_Found tests detection of the trial pack by name.
func TestHasTrialPack_Found(t *testing.T) {
	cr := &creditsSummary{
		Packages: []packageSummary{
			{Name: "CodeBuddy One-time Free 2-Week Pro Plan Trial", Remain: 250, Used: 0},
		},
	}
	if !hasTrialPack(cr) {
		t.Error("expected hasTrialPack=true for trial package")
	}
}

// TestHasTrialPack_NotFound tests that non-trial packages return false.
func TestHasTrialPack_NotFound(t *testing.T) {
	cr := &creditsSummary{
		Packages: []packageSummary{
			{Name: "CodeBuddy Free Monthly Pack", Remain: 500, Used: 0},
			{Name: "裂变奖励包A", Remain: 99, Used: 1},
		},
	}
	if hasTrialPack(cr) {
		t.Error("expected hasTrialPack=false for non-trial packages")
	}
}

// TestHasTrialPack_Nil tests nil safety.
func TestHasTrialPack_Nil(t *testing.T) {
	if hasTrialPack(nil) {
		t.Error("expected hasTrialPack=false for nil credits")
	}
}

// TestHasTrialPack_Empty tests empty packages.
func TestHasTrialPack_Empty(t *testing.T) {
	cr := &creditsSummary{}
	if hasTrialPack(cr) {
		t.Error("expected hasTrialPack=false for empty packages")
	}
}
