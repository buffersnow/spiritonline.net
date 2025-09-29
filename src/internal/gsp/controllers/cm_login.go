package controllers

import (
	"time"

	"buffersnow.com/spiritonline/internal/gsp/protocol"
	"buffersnow.com/spiritonline/pkg/gp"
)

func cm_Login(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error {

	if g.GPCM.LoggedIn {
		g.Log.Error("Login Error", "Authenticated user %d called login again", g.GPCM.AuthToken.WFCID)
		return protocol.GPLoginError_Duplicate
	}

	tokenPair := gci.Find("authtoken")
	if tokenPair.Length() == 0 {
		g.Log.Error("Login Error", "Length of authtoken kv was 0")
		return protocol.GPLoginError_InvalidToken
	}

	token, err := gp.DepickleWFCToken(tokenPair.Value().DataBlock())
	if err != nil {
		g.Log.Error("Login Error", "AuthToken is not in a valid format")
		return protocol.GPLoginError_InvalidToken
	}

	curTime := time.Now()
	if curTime.After(token.IssueTime.Add(1 * time.Minute)) {
		expiryTime := curTime.Sub(token.IssueTime) - (1 * time.Minute)

		g.Log.Error("Login Error", "AuthToken issued at %v expired %v ago", token.IssueTime, expiryTime)
		return protocol.GPLoginError_InvalidToken
	}

	return nil
}
