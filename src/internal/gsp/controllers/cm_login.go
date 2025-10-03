package controllers

import (
	"time"

	"buffersnow.com/spiritonline/internal/gsp/protocol"
	"buffersnow.com/spiritonline/pkg/gp"
)

func handleCM_Login(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error {

	if g.GPCM.LoggedIn {
		g.Log.Error("Login", "Authenticated user %d called login again", g.GPCM.AuthToken.WFCID)
		return protocol.GPLoginError_Duplicate
	}

	tokenPair := gci.Find("authtoken")
	if tokenPair.Length() == 0 {
		g.Log.Error("Login", "Length of authtoken kv was 0")
		return protocol.GPLoginError_InvalidToken
	}

	token, err := gp.DepickleWFCToken(tokenPair.Value().DataBlock())
	if err != nil {
		g.Log.Error("Login", "AuthToken is not in a valid format")
		return protocol.GPLoginError_InvalidToken
	}

	curTime := time.Now()
	expiryTime := token.IssueTime.Add(1 * time.Minute)

	if curTime.After(expiryTime) {
		timeSinceExpiry := curTime.Sub(expiryTime)

		g.Log.Error("Login", "AuthToken issued at %v expired %v ago", token.IssueTime, timeSinceExpiry)
		return protocol.GPLoginError_InvalidToken
	}

	g.GPCM.AuthToken = token

	return nil
}
