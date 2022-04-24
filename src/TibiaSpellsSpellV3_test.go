package main

import (
	"io"
	"testing"

	"github.com/TibiaData/tibiadata-api-go/src/static"
	"github.com/stretchr/testify/assert"
)

func TestFindPerson(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/spells/spell/Find Person.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	findPersonJson, err := TibiaSpellsSpellV3Impl("Find Person", string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal("", findPersonJson.Spells.Spell.Description)
	assert.Equal("Find Person", findPersonJson.Spells.Spell.Name)
	assert.Equal("find person", findPersonJson.Spells.Spell.Spell)
	assert.True(findPersonJson.Spells.Spell.HasSpellInformation)
	assert.NotNil(findPersonJson.Spells.Spell.SpellInformation)
	assert.Equal("exiva 'name'", findPersonJson.Spells.Spell.SpellInformation.Formula)
	assert.Equal(4, len(findPersonJson.Spells.Spell.SpellInformation.Vocation))
	assert.Equal("Druid", findPersonJson.Spells.Spell.SpellInformation.Vocation[0])
	assert.Equal("Knight", findPersonJson.Spells.Spell.SpellInformation.Vocation[1])
	assert.Equal("Paladin", findPersonJson.Spells.Spell.SpellInformation.Vocation[2])
	assert.Equal("Sorcerer", findPersonJson.Spells.Spell.SpellInformation.Vocation[3])
	assert.False(findPersonJson.Spells.Spell.SpellInformation.GroupAttack)
	assert.False(findPersonJson.Spells.Spell.SpellInformation.GroupHealing)
	assert.True(findPersonJson.Spells.Spell.SpellInformation.GroupSupport)
	assert.True(findPersonJson.Spells.Spell.SpellInformation.TypeInstant)
	assert.False(findPersonJson.Spells.Spell.SpellInformation.TypeRune)
	assert.Equal(2, findPersonJson.Spells.Spell.SpellInformation.CooldownAlone)
	assert.Equal(2, findPersonJson.Spells.Spell.SpellInformation.CooldownGroup)
	assert.Equal(8, findPersonJson.Spells.Spell.SpellInformation.Level)
	assert.Equal(20, findPersonJson.Spells.Spell.SpellInformation.Mana)
	assert.Equal(80, findPersonJson.Spells.Spell.SpellInformation.Price)
	assert.Equal(12, len(findPersonJson.Spells.Spell.SpellInformation.City))
	assert.Equal("Ab'Dendriel", findPersonJson.Spells.Spell.SpellInformation.City[0])
	assert.Equal("Yalahar", findPersonJson.Spells.Spell.SpellInformation.City[11])
	assert.False(findPersonJson.Spells.Spell.SpellInformation.Premium)
	assert.False(findPersonJson.Spells.Spell.HasRuneInformation)
}

func TestHeavyMagicMissileRune(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/spells/spell/Heavy Magic Missile Rune.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	hmmJson, err := TibiaSpellsSpellV3Impl("Heavy Magic Missile Rune", string(data))
	if err != nil {
		t.Fatal(err)
	}

	assert := assert.New(t)

	assert.Equal("", hmmJson.Spells.Spell.Description)
	assert.Equal("Heavy Magic Missile Rune", hmmJson.Spells.Spell.Name)
	assert.Equal("heavy magic missile rune", hmmJson.Spells.Spell.Spell)
	assert.True(hmmJson.Spells.Spell.HasSpellInformation)
	assert.NotNil(hmmJson.Spells.Spell.SpellInformation)
	assert.Equal("adori vis", hmmJson.Spells.Spell.SpellInformation.Formula)
	assert.Equal(2, len(hmmJson.Spells.Spell.SpellInformation.Vocation))
	assert.Equal("Druid", hmmJson.Spells.Spell.SpellInformation.Vocation[0])
	assert.Equal("Sorcerer", hmmJson.Spells.Spell.SpellInformation.Vocation[1])
	assert.False(hmmJson.Spells.Spell.SpellInformation.GroupAttack)
	assert.False(hmmJson.Spells.Spell.SpellInformation.GroupHealing)
	assert.True(hmmJson.Spells.Spell.SpellInformation.GroupSupport)
	assert.False(hmmJson.Spells.Spell.SpellInformation.TypeInstant)
	assert.True(hmmJson.Spells.Spell.SpellInformation.TypeRune)
	assert.Equal(2, hmmJson.Spells.Spell.SpellInformation.CooldownAlone)
	assert.Equal(2, hmmJson.Spells.Spell.SpellInformation.CooldownGroup)
	assert.Equal(2, hmmJson.Spells.Spell.SpellInformation.SoulPoints)
	assert.Equal(10, hmmJson.Spells.Spell.SpellInformation.Amount)
	assert.Equal(25, hmmJson.Spells.Spell.SpellInformation.Level)
	assert.Equal(350, hmmJson.Spells.Spell.SpellInformation.Mana)
	assert.Equal(1500, hmmJson.Spells.Spell.SpellInformation.Price)
	assert.Equal(13, len(hmmJson.Spells.Spell.SpellInformation.City))
	assert.Equal("Ab'Dendriel", hmmJson.Spells.Spell.SpellInformation.City[0])
	assert.Equal("Venore", hmmJson.Spells.Spell.SpellInformation.City[11])
	assert.False(hmmJson.Spells.Spell.SpellInformation.Premium)
	assert.True(hmmJson.Spells.Spell.HasRuneInformation)
	assert.Equal(4, len(hmmJson.Spells.Spell.RuneInformation.Vocation))
	assert.Equal("Druid", hmmJson.Spells.Spell.RuneInformation.Vocation[0])
	assert.Equal("Knight", hmmJson.Spells.Spell.RuneInformation.Vocation[1])
	assert.Equal("Paladin", hmmJson.Spells.Spell.RuneInformation.Vocation[2])
	assert.Equal("Sorcerer", hmmJson.Spells.Spell.RuneInformation.Vocation[3])
	assert.True(hmmJson.Spells.Spell.RuneInformation.GroupAttack)
	assert.False(hmmJson.Spells.Spell.RuneInformation.GroupHealing)
	assert.False(hmmJson.Spells.Spell.RuneInformation.GroupSupport)
	assert.Equal("energy", hmmJson.Spells.Spell.RuneInformation.DamageType)
	assert.Equal(25, hmmJson.Spells.Spell.RuneInformation.Level)
	assert.Equal(3, hmmJson.Spells.Spell.RuneInformation.MagicLevel)
}
