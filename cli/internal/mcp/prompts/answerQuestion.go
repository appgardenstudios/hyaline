package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// AnswerQuestionPrompt creates the answer_question prompt definition
func AnswerQuestionPrompt() mcp.Prompt {
	return mcp.NewPrompt("answer_question",
		mcp.WithPromptDescription("Answer a question based on available documentation"),
		mcp.WithArgument("question",
			mcp.ArgumentDescription("The question to answer"),
			mcp.RequiredArgument(),
		),
	)
}

// HandleAnswerQuestion handles the answer_question prompt
func HandleAnswerQuestion(_ context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	question := request.Params.Arguments["question"]
	if question == "" {
		return nil, fmt.Errorf("question is required")
	}

	return mcp.NewGetPromptResult(
		"Answer a question based on available documentation",
		[]mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(fmt.Sprintf("Based on the available documentation, please answer the following question: %s\n\nUse the list_documents and get_documents tools to find relevant information, then provide a comprehensive answer with links to source documentation.", question)),
			),
		},
	), nil
}
