// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build ignore_vet

package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/genai"
)

func run(ctx context.Context) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// RegisterFiles is only supported in the Gemini Developer client (mldev).
	if client.ClientConfig().Backend == genai.BackendVertexAI {
		log.Fatal("RegisterFiles is only supported in the Gemini Developer client.")
	}

	uris := []string{"gs://tensorflow_docs/image.jpg"}
	resp, err := client.Files.RegisterFiles(ctx, uris, nil)
	if err != nil {
		log.Fatal("Failed to register files: ", err)
	}

	fmt.Printf("Registered %d files.\n", len(resp.Files))
	if len(resp.Files) == 0 {
		log.Fatal("No files were registered.")
	}

	file := resp.Files[0]
	fmt.Printf("Registered file: %s\n", file.Name)

	// Use the registered file in a GenerateContent call.
	// We use the file's URI and MIMEType.
	parts := []*genai.Part{
		genai.NewPartFromText("can you summarize this file?"),
		genai.NewPartFromURI(file.URI, file.MIMEType),
	}
	content := genai.NewContentFromParts(parts, genai.RoleUser)

	fmt.Println("Generating content using the registered file...")
	genResp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", []*genai.Content{content}, nil)
	if err != nil {
		log.Fatal("Failed to generate content: ", err)
	}

	fmt.Printf("Response: %s\n", genResp.Text())
}

func main() {
	ctx := context.Background()
	run(ctx)
}
