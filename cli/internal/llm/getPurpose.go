package llm

import "fmt"

func GetDocumentPurpose(filename string, contents string) (string, error) {
	systemPrompt := "You are a senior technical writer who writes clear and accurate system documentation."
	prompt := fmt.Sprintf(`Write 1 to 2 sentences explaining the purpose of the document in <document>. The name of the document is %s. Be concise and accurate

<document>
%s
</document>`, filename, contents)

	fmt.Println(systemPrompt)
	fmt.Println(prompt)

	return "PURPOSE", nil
}

func GetSectionPurpose() (string, error) {
	return "PURPOSE", nil
}
