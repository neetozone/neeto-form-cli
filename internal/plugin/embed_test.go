package plugin

import (
	"os"
	"testing"
)

const repoSkillPath = "../../skills/neetoform/SKILL.md"

func TestRepoSkillMatchesEmbeddedSkill(t *testing.T) {
	repoSkill, err := os.ReadFile(repoSkillPath)
	if err != nil {
		t.Fatalf("could not read %s: %v", repoSkillPath, err)
	}

	if string(repoSkill) != skillContent {
		t.Fatalf(
			"%s and internal/plugin/skill.md have drifted.\n"+
				"skill.md is what `neetoform setup` installs into an assistant, so the two must "+
				"stay identical. Copy one over the other:\n"+
				"  cp internal/plugin/skill.md skills/neetoform/SKILL.md",
			repoSkillPath,
		)
	}
}
