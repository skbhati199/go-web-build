package framework

import (
	"context"
	"testing"
)

func TestReactAdapter(t *testing.T) {
	adapter, err := NewFrameworkAdapter("react")
	if err != nil {
		t.Fatalf("Failed to create React adapter: %v", err)
	}

	// Test initialization
	err = adapter.Initialize(context.Background(), Config{
		FrameworkType: "react",
		TemplateVersion: "1.0.0",
	})
	if err != nil {
		t.Errorf("Failed to initialize React adapter: %v", err)
	}

	// Test project generation
	err = adapter.GenerateProject(context.Background(), ProjectOptions{
		Name: "test-react-app",
		TypeScript: true,
	})
	if err != nil {
		t.Errorf("Failed to generate React project: %v", err)
	}

	// Test build process
	err = adapter.BuildProject(context.Background(), BuildModeDevelopment)
	if err != nil {
		t.Errorf("Failed to build React project: %v", err)
	}
}

func TestVueAdapter(t *testing.T) {
	adapter, err := NewFrameworkAdapter("vue")
	if err != nil {
		t.Fatalf("Failed to create Vue adapter: %v", err)
	}

	// Test initialization
	err = adapter.Initialize(context.Background(), Config{
		FrameworkType: "vue",
		TemplateVersion: "1.0.0",
	})
	if err != nil {
		t.Errorf("Failed to initialize Vue adapter: %v", err)
	}

	// Test project generation
	err = adapter.GenerateProject(context.Background(), ProjectOptions{
		Name: "test-vue-app",
		TypeScript: true,
	})
	if err != nil {
		t.Errorf("Failed to generate Vue project: %v", err)
	}

	// Test build process
	err = adapter.BuildProject(context.Background(), BuildModeDevelopment)
	if err != nil {
		t.Errorf("Failed to build Vue project: %v", err)
	}
}