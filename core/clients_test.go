package core

import "testing"

func TestAllClientsHasFive(t *testing.T) {
	want := []ClientKind{ClientBot, ClientAccount, ClientRubino, ClientTV, ClientRubiAI}
	if len(AllClients) != len(want) {
		t.Fatalf("AllClients len = %d, want %d", len(AllClients), len(want))
	}
	for i, k := range want {
		if AllClients[i].Kind != k {
			t.Errorf("AllClients[%d].Kind = %q, want %q", i, AllClients[i].Kind, k)
		}
	}
}

func TestEveryClientHasMetadata(t *testing.T) {
	for _, c := range AllClients {
		if c.Kind == "" || c.Name == "" || c.Package == "" || c.Description == "" || c.BaseURL == "" {
			t.Errorf("incomplete ClientInfo: %+v", c)
		}
	}
}

func TestClientFor(t *testing.T) {
	for _, c := range AllClients {
		got := ClientFor(c.Kind)
		if got.Kind != c.Kind || got.Name != c.Name || got.Package != c.Package {
			t.Errorf("ClientFor(%q) = %+v, want %+v", c.Kind, got, c)
		}
	}
	if got := ClientFor(ClientKind("nope")); got.Kind != "" {
		t.Errorf("ClientFor(unknown) = %+v, want zero value", got)
	}
}

func TestClientNamesAndKinds(t *testing.T) {
	names := ClientNames()
	kinds := ClientKinds()
	if len(names) != len(AllClients) || len(kinds) != len(AllClients) {
		t.Fatalf("names/kinds length mismatch: %d / %d", len(names), len(kinds))
	}
	for i, c := range AllClients {
		if names[i] != c.Name {
			t.Errorf("ClientNames()[%d] = %q, want %q", i, names[i], c.Name)
		}
		if kinds[i] != c.Kind {
			t.Errorf("ClientKinds()[%d] = %q, want %q", i, kinds[i], c.Kind)
		}
	}
}

func TestClientsAreRubikaOnly(t *testing.T) {
	// Every first-class client must talk to the Rubika ecosystem.
	for _, c := range AllClients {
		if c.BaseURL == "" {
			t.Errorf("ClientInfo %q missing BaseURL", c.Kind)
		}
	}
}
