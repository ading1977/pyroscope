package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	otelProfile "go.opentelemetry.io/proto/otlp/profiles/v1development"
	"google.golang.org/protobuf/proto"
)

func debugOTLPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	var request otelProfile.ExportProfilesServiceRequest
	if err := proto.Unmarshal(body, &request); err != nil {
		http.Error(w, "Failed to unmarshal OTLP", http.StatusBadRequest)
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("\n=== OTLP DEBUG [%s] ===\n", timestamp)
	
	totalMappings := 0
	unsymbolizedMappings := 0
	
	for i, resourceProfile := range request.ResourceProfiles {
		fmt.Printf("ResourceProfile %d:\n", i)
		for j, scopeProfile := range resourceProfile.ScopeProfiles {
			fmt.Printf("  ScopeProfile %d:\n", j)
			for k, profile := range scopeProfile.Profiles {
				fmt.Printf("    Profile %d:\n", k)
				dict := profile.Profile
				
				// Check mappings
				for l, mapping := range dict.MappingTable {
					totalMappings++
					filename := "unknown"
					if mapping.FilenameStrindex < int32(len(dict.StringTable)) {
						filename = dict.StringTable[mapping.FilenameStrindex]
					}
					
					fmt.Printf("      Mapping %d: HasFunctions=%v, HasFilenames=%v, HasLineNumbers=%v\n", 
						l, mapping.HasFunctions, mapping.HasFilenames, mapping.HasLineNumbers)
					fmt.Printf("        MemoryStart=0x%x, MemoryLimit=0x%x, Filename=%s\n", 
						mapping.MemoryStart, mapping.MemoryLimit, filename)
					
					// Check attributes for build ID
					buildID := "none"
					for _, attrIdx := range mapping.AttributeIndices {
						if attrIdx < int32(len(dict.AttributeTable)) {
							attr := dict.AttributeTable[attrIdx]
							if attr.Key == "process.executable.build_id.gnu" {
								buildID = attr.Value.GetStringValue()
							}
						}
					}
					fmt.Printf("        BuildID=%s\n", buildID)
					
					if !mapping.HasFunctions {
						unsymbolizedMappings++
					}
				}
				
				// Check locations
				fmt.Printf("      Locations: %d\n", len(dict.LocationTable))
				unsymbolizedLocations := 0
				for l, location := range dict.LocationTable {
					if len(location.Line) == 0 {
						unsymbolizedLocations++
					}
					if l < 5 { // Show first 5 for debugging
						fmt.Printf("        Location %d: Address=0x%x, MappingIndex=%d, Lines=%d\n", 
							l, location.Address, location.MappingIndex, len(location.Line))
					}
				}
				fmt.Printf("      Unsymbolized locations: %d/%d\n", unsymbolizedLocations, len(dict.LocationTable))
			}
		}
	}
	
	fmt.Printf("SUMMARY: Total mappings=%d, Unsymbolized mappings=%d\n", totalMappings, unsymbolizedMappings)
	fmt.Printf("=== END DEBUG ===\n\n")
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Debug complete - check logs"))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	fmt.Println("Starting OTLP Debug Server on :4040")
	fmt.Println("POST profiles to /v1/profiles for debugging")
	fmt.Println("Health check at /health")
	
	http.HandleFunc("/v1/profiles", debugOTLPHandler)
	http.HandleFunc("/health", healthHandler)
	
	log.Fatal(http.ListenAndServe(":4040", nil))
}