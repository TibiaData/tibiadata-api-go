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

	assert.Empty(findPersonJson.Spells.Spell.Description)
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

	assert.Empty(hmmJson.Spells.Spell.Description)
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

func TestAnnihilation(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/spells/spell/Annihilation.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	annihilationJson, _ := TibiaSpellsSpellV3Impl("Annihilation", string(data))
	assert := assert.New(t)

	assert.Empty(annihilationJson.Spells.Spell.Description)
	assert.Equal("Annihilation", annihilationJson.Spells.Spell.Name)
	assert.Equal("annihilation", annihilationJson.Spells.Spell.Spell)
	assert.True(annihilationJson.Spells.Spell.HasSpellInformation)
	assert.NotNil(annihilationJson.Spells.Spell.SpellInformation)
	assert.Equal("exori gran ico", annihilationJson.Spells.Spell.SpellInformation.Formula)
	assert.Equal(1, len(annihilationJson.Spells.Spell.SpellInformation.Vocation))
	assert.Equal("Knight", annihilationJson.Spells.Spell.SpellInformation.Vocation[0])
	assert.True(annihilationJson.Spells.Spell.SpellInformation.GroupAttack)
	assert.False(annihilationJson.Spells.Spell.SpellInformation.GroupHealing)
	assert.False(annihilationJson.Spells.Spell.SpellInformation.GroupSupport)
	assert.True(annihilationJson.Spells.Spell.SpellInformation.TypeInstant)
	assert.False(annihilationJson.Spells.Spell.SpellInformation.TypeRune)
	assert.Equal("var.", annihilationJson.Spells.Spell.SpellInformation.DamageType) // weird one..
	assert.Equal(30, annihilationJson.Spells.Spell.SpellInformation.CooldownAlone)
	assert.Equal(4, annihilationJson.Spells.Spell.SpellInformation.CooldownGroup)
	assert.Equal(0, annihilationJson.Spells.Spell.SpellInformation.SoulPoints)
	assert.Equal(0, annihilationJson.Spells.Spell.SpellInformation.Amount)
	assert.Equal(7, len(annihilationJson.Spells.Spell.SpellInformation.City))
	assert.True(annihilationJson.Spells.Spell.SpellInformation.Premium)
	assert.False(annihilationJson.Spells.Spell.HasRuneInformation)
}

func TestBruiseBane(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/spells/spell/Bruise Bane.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	bruisebaneJson, _ := TibiaSpellsSpellV3Impl("Bruise Bane", string(data))
	assert := assert.New(t)

	assert.Empty(bruisebaneJson.Spells.Spell.Description)
	assert.Equal("Bruise Bane", bruisebaneJson.Spells.Spell.Name)
	assert.Equal("bruise bane", bruisebaneJson.Spells.Spell.Spell)
	assert.True(bruisebaneJson.Spells.Spell.HasSpellInformation)
	assert.NotNil(bruisebaneJson.Spells.Spell.SpellInformation)
	assert.Equal("exura infir ico", bruisebaneJson.Spells.Spell.SpellInformation.Formula)
	assert.Equal(1, len(bruisebaneJson.Spells.Spell.SpellInformation.Vocation))
	assert.False(bruisebaneJson.Spells.Spell.SpellInformation.GroupAttack)
	assert.True(bruisebaneJson.Spells.Spell.SpellInformation.GroupHealing)
	assert.False(bruisebaneJson.Spells.Spell.SpellInformation.GroupSupport)
	assert.True(bruisebaneJson.Spells.Spell.SpellInformation.TypeInstant)
	assert.False(bruisebaneJson.Spells.Spell.SpellInformation.TypeRune)
	assert.Equal(1, bruisebaneJson.Spells.Spell.SpellInformation.CooldownAlone)
	assert.Equal(1, bruisebaneJson.Spells.Spell.SpellInformation.CooldownGroup)
	assert.Equal(1, bruisebaneJson.Spells.Spell.SpellInformation.Level)
	assert.Equal(10, bruisebaneJson.Spells.Spell.SpellInformation.Mana)
	assert.Equal(0, bruisebaneJson.Spells.Spell.SpellInformation.Price)
	assert.Equal(1, len(bruisebaneJson.Spells.Spell.SpellInformation.City))
	assert.Equal("Dawnport", bruisebaneJson.Spells.Spell.SpellInformation.City[0])
	assert.False(bruisebaneJson.Spells.Spell.SpellInformation.Premium)
	assert.False(bruisebaneJson.Spells.Spell.HasRuneInformation)
}

func TestCurePoisonRune(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/spells/spell/Cure Poison Rune.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	curepoisonruneJson, _ := TibiaSpellsSpellV3Impl("Cure Poison Rune", string(data))
	assert := assert.New(t)

	assert.Empty(curepoisonruneJson.Spells.Spell.Description)
	assert.Equal("Cure Poison Rune", curepoisonruneJson.Spells.Spell.Name)
	assert.Equal("cure poison rune", curepoisonruneJson.Spells.Spell.Spell)
	assert.True(curepoisonruneJson.Spells.Spell.HasSpellInformation)
	assert.NotNil(curepoisonruneJson.Spells.Spell.SpellInformation)
	assert.Equal("adana pox", curepoisonruneJson.Spells.Spell.SpellInformation.Formula)
	assert.Equal(1, len(curepoisonruneJson.Spells.Spell.SpellInformation.Vocation))
	assert.False(curepoisonruneJson.Spells.Spell.SpellInformation.GroupAttack)
	assert.False(curepoisonruneJson.Spells.Spell.SpellInformation.GroupHealing)
	assert.True(curepoisonruneJson.Spells.Spell.SpellInformation.GroupSupport)
	assert.False(curepoisonruneJson.Spells.Spell.SpellInformation.TypeInstant)
	assert.True(curepoisonruneJson.Spells.Spell.SpellInformation.TypeRune)
	assert.False(curepoisonruneJson.Spells.Spell.SpellInformation.Premium)
	assert.True(curepoisonruneJson.Spells.Spell.HasRuneInformation)
	assert.False(curepoisonruneJson.Spells.Spell.RuneInformation.GroupAttack)
	assert.True(curepoisonruneJson.Spells.Spell.RuneInformation.GroupHealing)
	assert.False(curepoisonruneJson.Spells.Spell.RuneInformation.GroupSupport)
}

func TestConvinceCreatureRune(t *testing.T) {
	file, err := static.TestFiles.Open("testdata/spells/spell/Convince Creature Rune.html")
	if err != nil {
		t.Fatalf("file opening error: %s", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("File reading error: %s", err)
	}

	convincecreatureruneJson, _ := TibiaSpellsSpellV3Impl("Convince Creature Rune", string(data))
	assert := assert.New(t)

	assert.Empty(convincecreatureruneJson.Spells.Spell.Description)
	assert.Equal("Convince Creature Rune", convincecreatureruneJson.Spells.Spell.Name)
	assert.Equal("convince creature rune", convincecreatureruneJson.Spells.Spell.Spell)
	assert.True(convincecreatureruneJson.Spells.Spell.HasSpellInformation)
	assert.NotNil(convincecreatureruneJson.Spells.Spell.SpellInformation)
	assert.Equal("adeta sio", convincecreatureruneJson.Spells.Spell.SpellInformation.Formula)
	assert.True(convincecreatureruneJson.Spells.Spell.HasRuneInformation)
	assert.False(convincecreatureruneJson.Spells.Spell.RuneInformation.GroupAttack)
	assert.False(convincecreatureruneJson.Spells.Spell.RuneInformation.GroupHealing)
	assert.True(convincecreatureruneJson.Spells.Spell.RuneInformation.GroupSupport)
}
