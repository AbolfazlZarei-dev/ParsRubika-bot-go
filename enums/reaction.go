package enums

// Reaction is the type of message reaction. The numeric ID (used on the wire) is exposed via ReactionID constants and
// the ReactionFromID lookup.
type Reaction string

const (
	ReactionNone    Reaction = ""
	ReactionLike    Reaction = "Like"
	ReactionDislike Reaction = "Dislike"
	ReactionHeart   Reaction = "Heart"
	ReactionFire    Reaction = "Fire"
)

// String implements fmt.Stringer.
func (r Reaction) String() string { return string(r) }

// ReactionID is the numeric reaction identifier used on the wire. The table
// covers all 72 reactions. Use ReactionFromID to resolve an ID to its name.
type ReactionID int

const (
	ReactionIDRedHeart              ReactionID = 1
	ReactionIDThumbsUp              ReactionID = 2 // Like
	ReactionIDThumbsDown            ReactionID = 3 // Dislike / UnLike
	ReactionIDFire                  ReactionID = 4
	ReactionIDSmilingFaceWithHearts ReactionID = 5
	ReactionIDClappingHands         ReactionID = 6
	ReactionIDBeamingFace           ReactionID = 7
	ReactionIDThinkingFace          ReactionID = 8
	ReactionIDFaceWithTears         ReactionID = 9
	ReactionIDRollingFace           ReactionID = 10
	ReactionIDWinkingFace           ReactionID = 11
	ReactionIDHappyFace             ReactionID = 12
	ReactionIDSmirkingFace          ReactionID = 13
	ReactionIDMeltingFace           ReactionID = 14
	ReactionIDNauseatedFace         ReactionID = 15
	ReactionIDMoneyFace             ReactionID = 16
	ReactionIDPartyingFace          ReactionID = 17
	ReactionIDConfoundedFace        ReactionID = 18
	ReactionIDHundredPoints         ReactionID = 19
	ReactionIDHuggingFace           ReactionID = 20
	ReactionIDSalutingFace          ReactionID = 21
	ReactionIDFlexedBiceps          ReactionID = 22
	ReactionIDFoldedHands           ReactionID = 23
	ReactionIDOKHand                ReactionID = 24
	ReactionIDRaisedFist            ReactionID = 25
	ReactionIDPeaceHand             ReactionID = 26
	ReactionIDWaveHand              ReactionID = 27
	ReactionIDBrain                 ReactionID = 28
	ReactionIDRedHeart2             ReactionID = 29
	ReactionIDGrowingHeart          ReactionID = 30
	ReactionIDBrokenHeart           ReactionID = 31
	ReactionIDSplashingSweat        ReactionID = 32
	ReactionIDRocket                ReactionID = 33
	ReactionIDCollision             ReactionID = 34
	ReactionIDGlowingStar           ReactionID = 35
	ReactionIDPartyPopper           ReactionID = 36
	ReactionIDTrophy                ReactionID = 37
	ReactionIDFirstPlace            ReactionID = 38
	ReactionIDMedal                 ReactionID = 39
	ReactionIDBullseye              ReactionID = 40
	ReactionIDDart                  ReactionID = 41
	ReactionIDCheckMarkButton       ReactionID = 42
	ReactionIDCrossMark             ReactionID = 43
	ReactionIDOKButton              ReactionID = 44
	ReactionIDSnowflake             ReactionID = 45
	ReactionIDJackOLantern          ReactionID = 46
	ReactionIDFirecracker           ReactionID = 47
	ReactionIDPileOfPoo             ReactionID = 48
	ReactionIDVolcano               ReactionID = 49
	ReactionIDDog                   ReactionID = 50
	ReactionIDCat                   ReactionID = 51
	ReactionIDUnicorn               ReactionID = 52
	ReactionIDPanda                 ReactionID = 53
	ReactionIDRabbit                ReactionID = 54
	ReactionIDPig                   ReactionID = 55
	ReactionIDChicken               ReactionID = 56
	ReactionIDDuck                  ReactionID = 57
	ReactionIDBug                   ReactionID = 58
	ReactionIDFish                  ReactionID = 59
	ReactionIDWhale                 ReactionID = 60
	ReactionIDDolphin               ReactionID = 61
	ReactionIDRose                  ReactionID = 62
	ReactionIDTulip                 ReactionID = 63
	ReactionIDSunflower             ReactionID = 64
	ReactionIDCactus                ReactionID = 65
	ReactionIDBouquet               ReactionID = 66
	ReactionIDWiltedFlower          ReactionID = 67
	ReactionIDMapleLeaf             ReactionID = 68
	ReactionIDFourLeafClover        ReactionID = 69
	ReactionIDHotBeverage           ReactionID = 70
	ReactionIDBirthdayCake          ReactionID = 71
	ReactionIDRage                  ReactionID = 72
)

// reactionIDNames maps every numeric reaction ID to a human-readable name.
// This is a package-level table so ReactionFromID is a fast lookup.
var reactionIDNames = map[ReactionID]string{
	ReactionIDRedHeart:              "Red Heart",
	ReactionIDThumbsUp:              "Thumbs Up",
	ReactionIDThumbsDown:            "Thumbs Down",
	ReactionIDFire:                  "Fire",
	ReactionIDSmilingFaceWithHearts: "Smiling Face with Hearts",
	ReactionIDClappingHands:         "Clapping Hands",
	ReactionIDBeamingFace:           "Beaming Face",
	ReactionIDThinkingFace:          "Thinking Face",
	ReactionIDFaceWithTears:         "Face with Tears",
	ReactionIDRollingFace:           "Rolling Face",
	ReactionIDWinkingFace:           "Winking Face",
	ReactionIDHappyFace:             "Happy Face",
	ReactionIDSmirkingFace:          "Smirking Face",
	ReactionIDMeltingFace:           "Melting Face",
	ReactionIDNauseatedFace:         "Nauseated Face",
	ReactionIDMoneyFace:             "Money Face",
	ReactionIDPartyingFace:          "Partying Face",
	ReactionIDConfoundedFace:        "Confounded Face",
	ReactionIDHundredPoints:         "Hundred Points",
	ReactionIDHuggingFace:           "Hugging Face",
	ReactionIDSalutingFace:          "Saluting Face",
	ReactionIDFlexedBiceps:          "Flexed Biceps",
	ReactionIDFoldedHands:           "Folded Hands",
	ReactionIDOKHand:                "OK Hand",
	ReactionIDRaisedFist:            "Raised Fist",
	ReactionIDPeaceHand:             "Peace Hand",
	ReactionIDWaveHand:              "Wave Hand",
	ReactionIDBrain:                 "Brain",
	ReactionIDRedHeart2:             "Red Heart",
	ReactionIDGrowingHeart:          "Growing Heart",
	ReactionIDBrokenHeart:           "Broken Heart",
	ReactionIDSplashingSweat:        "Splashing Sweat",
	ReactionIDRocket:                "Rocket",
	ReactionIDCollision:             "Collision",
	ReactionIDGlowingStar:           "Glowing Star",
	ReactionIDPartyPopper:           "Party Popper",
	ReactionIDTrophy:                "Trophy",
	ReactionIDFirstPlace:            "First Place",
	ReactionIDMedal:                 "Medal",
	ReactionIDBullseye:              "Bullseye",
	ReactionIDDart:                  "Dart",
	ReactionIDCheckMarkButton:       "Check Mark Button",
	ReactionIDCrossMark:             "Cross Mark",
	ReactionIDOKButton:              "OK Button",
	ReactionIDSnowflake:             "Snowflake",
	ReactionIDJackOLantern:          "Jack-O-Lantern",
	ReactionIDFirecracker:           "Firecracker",
	ReactionIDPileOfPoo:             "Pile of Poo",
	ReactionIDVolcano:               "Volcano",
	ReactionIDDog:                   "Dog",
	ReactionIDCat:                   "Cat",
	ReactionIDUnicorn:               "Unicorn",
	ReactionIDPanda:                 "Panda",
	ReactionIDRabbit:                "Rabbit",
	ReactionIDPig:                   "Pig",
	ReactionIDChicken:               "Chicken",
	ReactionIDDuck:                  "Duck",
	ReactionIDBug:                   "Bug",
	ReactionIDFish:                  "Fish",
	ReactionIDWhale:                 "Whale",
	ReactionIDDolphin:               "Dolphin",
	ReactionIDRose:                  "Rose",
	ReactionIDTulip:                 "Tulip",
	ReactionIDSunflower:             "Sunflower",
	ReactionIDCactus:                "Cactus",
	ReactionIDBouquet:               "Bouquet",
	ReactionIDWiltedFlower:          "Wilted Flower",
	ReactionIDMapleLeaf:             "Maple Leaf",
	ReactionIDFourLeafClover:        "Four Leaf Clover",
	ReactionIDHotBeverage:           "Hot Beverage",
	ReactionIDBirthdayCake:          "Birthday Cake",
	ReactionIDRage:                  "Rage",
}

// String returns the human-readable name of the reaction.
func (r ReactionID) String() string {
	if name, ok := reactionIDNames[r]; ok {
		return name
	}
	return "Unknown"
}

// ReactionFromID resolves a numeric reaction ID to its name. Returns
// "Unknown" for unrecognized IDs.
func ReactionFromID(id int) string {
	return ReactionID(id).String()
}
