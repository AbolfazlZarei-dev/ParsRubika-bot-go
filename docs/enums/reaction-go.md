# `enums` — Reaction

_Source: [`reaction.go`](../../ParsRubika/enums/reaction.go)_

## Symbols

### `Reaction`

```go
type Reaction string
```

Reaction is the type of message reaction. The numeric ID (used on the wire) is exposed via ReactionID constants and
the ReactionFromID lookup.

### `ReactionNone`

```go
const ReactionNone Reaction = ""
```

### `ReactionLike`

```go
const ReactionLike Reaction = "Like"
```

### `ReactionDislike`

```go
const ReactionDislike Reaction = "Dislike"
```

### `ReactionHeart`

```go
const ReactionHeart Reaction = "Heart"
```

### `ReactionFire`

```go
const ReactionFire Reaction = "Fire"
```

### `String`

```go
func (r Reaction) String() string { return string(r) }
```

String implements fmt.Stringer.

### `ReactionID`

```go
type ReactionID int
```

ReactionID is the numeric reaction identifier used on the wire. The table
covers all 72 reactions. Use ReactionFromID to resolve an ID to its name.

### `ReactionIDRedHeart`

```go
const ReactionIDRedHeart ReactionID = 1
```

### `ReactionIDThumbsUp`

```go
const ReactionIDThumbsUp ReactionID = 2 // Like
```

### `ReactionIDThumbsDown`

```go
const ReactionIDThumbsDown ReactionID = 3 // Dislike / UnLike
```

### `ReactionIDFire`

```go
const ReactionIDFire ReactionID = 4
```

### `ReactionIDSmilingFaceWithHearts`

```go
const ReactionIDSmilingFaceWithHearts ReactionID = 5
```

### `ReactionIDClappingHands`

```go
const ReactionIDClappingHands ReactionID = 6
```

### `ReactionIDBeamingFace`

```go
const ReactionIDBeamingFace ReactionID = 7
```

### `ReactionIDThinkingFace`

```go
const ReactionIDThinkingFace ReactionID = 8
```

### `ReactionIDFaceWithTears`

```go
const ReactionIDFaceWithTears ReactionID = 9
```

### `ReactionIDRollingFace`

```go
const ReactionIDRollingFace ReactionID = 10
```

### `ReactionIDWinkingFace`

```go
const ReactionIDWinkingFace ReactionID = 11
```

### `ReactionIDHappyFace`

```go
const ReactionIDHappyFace ReactionID = 12
```

### `ReactionIDSmirkingFace`

```go
const ReactionIDSmirkingFace ReactionID = 13
```

### `ReactionIDMeltingFace`

```go
const ReactionIDMeltingFace ReactionID = 14
```

### `ReactionIDNauseatedFace`

```go
const ReactionIDNauseatedFace ReactionID = 15
```

### `ReactionIDMoneyFace`

```go
const ReactionIDMoneyFace ReactionID = 16
```

### `ReactionIDPartyingFace`

```go
const ReactionIDPartyingFace ReactionID = 17
```

### `ReactionIDConfoundedFace`

```go
const ReactionIDConfoundedFace ReactionID = 18
```

### `ReactionIDHundredPoints`

```go
const ReactionIDHundredPoints ReactionID = 19
```

### `ReactionIDHuggingFace`

```go
const ReactionIDHuggingFace ReactionID = 20
```

### `ReactionIDSalutingFace`

```go
const ReactionIDSalutingFace ReactionID = 21
```

### `ReactionIDFlexedBiceps`

```go
const ReactionIDFlexedBiceps ReactionID = 22
```

### `ReactionIDFoldedHands`

```go
const ReactionIDFoldedHands ReactionID = 23
```

### `ReactionIDOKHand`

```go
const ReactionIDOKHand ReactionID = 24
```

### `ReactionIDRaisedFist`

```go
const ReactionIDRaisedFist ReactionID = 25
```

### `ReactionIDPeaceHand`

```go
const ReactionIDPeaceHand ReactionID = 26
```

### `ReactionIDWaveHand`

```go
const ReactionIDWaveHand ReactionID = 27
```

### `ReactionIDBrain`

```go
const ReactionIDBrain ReactionID = 28
```

### `ReactionIDRedHeart2`

```go
const ReactionIDRedHeart2 ReactionID = 29
```

### `ReactionIDGrowingHeart`

```go
const ReactionIDGrowingHeart ReactionID = 30
```

### `ReactionIDBrokenHeart`

```go
const ReactionIDBrokenHeart ReactionID = 31
```

### `ReactionIDSplashingSweat`

```go
const ReactionIDSplashingSweat ReactionID = 32
```

### `ReactionIDRocket`

```go
const ReactionIDRocket ReactionID = 33
```

### `ReactionIDCollision`

```go
const ReactionIDCollision ReactionID = 34
```

### `ReactionIDGlowingStar`

```go
const ReactionIDGlowingStar ReactionID = 35
```

### `ReactionIDPartyPopper`

```go
const ReactionIDPartyPopper ReactionID = 36
```

### `ReactionIDTrophy`

```go
const ReactionIDTrophy ReactionID = 37
```

### `ReactionIDFirstPlace`

```go
const ReactionIDFirstPlace ReactionID = 38
```

### `ReactionIDMedal`

```go
const ReactionIDMedal ReactionID = 39
```

### `ReactionIDBullseye`

```go
const ReactionIDBullseye ReactionID = 40
```

### `ReactionIDDart`

```go
const ReactionIDDart ReactionID = 41
```

### `ReactionIDCheckMarkButton`

```go
const ReactionIDCheckMarkButton ReactionID = 42
```

### `ReactionIDCrossMark`

```go
const ReactionIDCrossMark ReactionID = 43
```

### `ReactionIDOKButton`

```go
const ReactionIDOKButton ReactionID = 44
```

### `ReactionIDSnowflake`

```go
const ReactionIDSnowflake ReactionID = 45
```

### `ReactionIDJackOLantern`

```go
const ReactionIDJackOLantern ReactionID = 46
```

### `ReactionIDFirecracker`

```go
const ReactionIDFirecracker ReactionID = 47
```

### `ReactionIDPileOfPoo`

```go
const ReactionIDPileOfPoo ReactionID = 48
```

### `ReactionIDVolcano`

```go
const ReactionIDVolcano ReactionID = 49
```

### `ReactionIDDog`

```go
const ReactionIDDog ReactionID = 50
```

### `ReactionIDCat`

```go
const ReactionIDCat ReactionID = 51
```

### `ReactionIDUnicorn`

```go
const ReactionIDUnicorn ReactionID = 52
```

### `ReactionIDPanda`

```go
const ReactionIDPanda ReactionID = 53
```

### `ReactionIDRabbit`

```go
const ReactionIDRabbit ReactionID = 54
```

### `ReactionIDPig`

```go
const ReactionIDPig ReactionID = 55
```

### `ReactionIDChicken`

```go
const ReactionIDChicken ReactionID = 56
```

### `ReactionIDDuck`

```go
const ReactionIDDuck ReactionID = 57
```

### `ReactionIDBug`

```go
const ReactionIDBug ReactionID = 58
```

### `ReactionIDFish`

```go
const ReactionIDFish ReactionID = 59
```

### `ReactionIDWhale`

```go
const ReactionIDWhale ReactionID = 60
```

### `ReactionIDDolphin`

```go
const ReactionIDDolphin ReactionID = 61
```

### `ReactionIDRose`

```go
const ReactionIDRose ReactionID = 62
```

### `ReactionIDTulip`

```go
const ReactionIDTulip ReactionID = 63
```

### `ReactionIDSunflower`

```go
const ReactionIDSunflower ReactionID = 64
```

### `ReactionIDCactus`

```go
const ReactionIDCactus ReactionID = 65
```

### `ReactionIDBouquet`

```go
const ReactionIDBouquet ReactionID = 66
```

### `ReactionIDWiltedFlower`

```go
const ReactionIDWiltedFlower ReactionID = 67
```

### `ReactionIDMapleLeaf`

```go
const ReactionIDMapleLeaf ReactionID = 68
```

### `ReactionIDFourLeafClover`

```go
const ReactionIDFourLeafClover ReactionID = 69
```

### `ReactionIDHotBeverage`

```go
const ReactionIDHotBeverage ReactionID = 70
```

### `ReactionIDBirthdayCake`

```go
const ReactionIDBirthdayCake ReactionID = 71
```

### `ReactionIDRage`

```go
const ReactionIDRage ReactionID = 72
```

### `String`

```go
func (r ReactionID) String() string
```

String returns the human-readable name of the reaction.

### `ReactionFromID`

```go
func ReactionFromID(id int) string
```

ReactionFromID resolves a numeric reaction ID to its name. Returns
"Unknown" for unrecognized IDs.
