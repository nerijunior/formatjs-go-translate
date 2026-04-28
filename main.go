package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"flag"
	"os"

	"strings"
	
	"github.com/aquasecurity/table"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
) 

type LocaleObject struct {
    Context string `json:"context"`
    String  string `json:"string"`
}

func main () {
	// Flags
	sourcePath := flag.String("source-file", "locales/defaultMessages.json", "defaultMessages json file")
	targetPath := flag.String("locale-file", "locales/pt-BR.json", "destination locale json file")
	diffOnly := flag.Bool("diff", false, "Show a diff only")

	flag.Parse()

	fmt.Printf("Comparing translations between %s and %s \n", *sourcePath, *targetPath)

	// Table
	table := table.New(os.Stdout)

	table.SetHeaders("KEY", "source", "translation")
	
	// Source files
	sourceData := loadJSON(*sourcePath)
	targetData := loadJSON(*targetPath)

	if len(sourceData) == 0 {
		fmt.Println("Erro: Arquivo de origem vazio ou não encontrado.")
		return
	}

	if *diffOnly {
		for key, _ := range sourceData {
			table.AddRow(key, sourceData[key].String, targetData[key].String)
		}
		table.Render()
		return
	}

	fmt.Println("Verificando chaves para tradução...")

	ctx := context.Background()

	credsFile := "google-creds.json"
	client, err := translate.NewClient(ctx, option.WithCredentialsFile(credsFile))
	if err != nil {
		log.Fatalf("Erro ao inicializar cliente: %v", err)
	}

	defer client.Close()

	fmt.Println("Cliente autenticado")
	hasChanges := translateInBatches(ctx, client, sourceData, targetData)


	if hasChanges {
		saveTargetJSON(*targetPath, targetData)
		fmt.Println("Arquivo pt-BR.json atualizado com sucesso!")
	} else {
		fmt.Println("Nenhuma chave nova encontrada")
	}

}

var batchSize = 100

func translateInBatches(ctx context.Context, client *translate.Client, source map[string]*LocaleObject, target map[string]*LocaleObject) bool {
	var keysToTranslate []string
	var valuesToTranslate []string
	hasChanges := false

	for key, unit := range source {
		existingTranslation, exists := target[key]
		isTodo := exists && strings.HasPrefix(existingTranslation.String, "[TODO]")

		if !exists || isTodo {
			keysToTranslate = append(keysToTranslate, key)
			valuesToTranslate = append(valuesToTranslate, unit.String)
		}

		if len(valuesToTranslate) >= batchSize {
			performTranslation(ctx, client, keysToTranslate, valuesToTranslate, target)
			keysToTranslate = nil
			valuesToTranslate = nil
			hasChanges = true
		}
	}

	if len(valuesToTranslate) > 0 {
		performTranslation(ctx, client, keysToTranslate, valuesToTranslate, target)
		hasChanges = true
	}

	return hasChanges
}

func performTranslation(ctx context.Context, client *translate.Client, keys []string, values []string, target map[string]*LocaleObject) {
	fmt.Printf("Enviando lote de %s strings para o Google...\n", len(values))

	resp, err := client.Translate(ctx, values, language.BrazilianPortuguese, nil)
	if err != nil {
		log.Printf("Erro ao traduzir lote: %v", err)
		return
	}

	for i, translation := range resp {
		key := keys[i]
		if target[key] == nil {
			target[key] = &LocaleObject{}
		}
		target[key].String = translation.Text
	}
}

func loadJSON(path string) map[string]*LocaleObject {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return make(map[string]*LocaleObject)
	}

	var data map[string]*LocaleObject
	json.Unmarshal(content, &data)
	
	return data
}

func saveTargetJSON(path string, data map[string]*LocaleObject){
	content, _ := json.MarshalIndent(data, "", " ")
	ioutil.WriteFile(path, content, 0644)
}
