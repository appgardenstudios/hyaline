package llm

import (
	"fmt"
	"hyaline/internal/config"
)

func GetDocumentPurpose(filename string, contents string, cfg *config.LLM) (result string, err error) {
	systemPrompt := "You are a senior technical writer who writes clear and accurate system documentation."
	userPrompt := fmt.Sprintf(`Write 1 to 2 sentences explaining the purpose of the document in <document>. The name of the document is %s. Be concise and accurate

<document>
%s
</document>`, filename, contents)

	result, err = callLLM(systemPrompt, userPrompt, cfg)

	return
}

func GetSectionPurpose() (string, error) {
	return "PURPOSE", nil
}
