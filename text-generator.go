package main

import (
	"math/rand/v2"
)

var quotes = [...]string{
	"The gentle breeze carried the scent of blooming jasmine through the open window, filling the room with a calming aroma that made it hard not to feel at peace with the world outside. As the sunlight streamed in, dancing across the wooden floor, she found herself lost in thought, reminiscing about simpler days spent wandering through lush fields, where the same scent of jasmine seemed to follow her every step.",
	"A mysterious figure appeared at the edge of the forest, shrouded in a thick mist that obscured all but the faint outline of what looked like a long coat and a wide-brimmed hat. The air around them seemed unnaturally still, as if the presence of this stranger had paused time itself, leaving only the faint sound of rustling leaves to hint at the unseen movements deeper within the woods",
	"She couldn't believe her luck when she found a rare coin buried in the sand, its intricate designs glinting under the afternoon sun as if it were trying to tell her a story from centuries ago. The discovery brought a rush of excitement, and she began to imagine the hands that might have held this coin before her: merchants, explorers, perhaps even royalty, all connected by this small piece of history now resting in her palm!!!",
	"The cat perched on the windowsill, watching the world go by with unblinking eyes, as though it were silently judging the hurried lives of the people passing below on the busy street. Every now and then, it flicked its tail with an air of regal disdain, as if to remind anyone who cared to notice that it had no interest in the trivial concerns of the bustling world beyond the glass.",
	"The old library was a treasure trove of forgotten knowledge, with towering shelves of dusty books that seemed to whisper secrets to anyone curious enough to open them. The faint smell of aged paper and leather bindings created an atmosphere of timeless wonder, as though stepping inside transported you to a realm where stories came alive and the boundaries of imagination were limited only by the words on the pages.",
	"As the thunderstorm raged outside, the rhythmic drumming of rain on the roof provided an oddly soothing background to the crackling fire that warmed the small cabin. Lightning illuminated the darkened room at intervals, casting fleeting shadows on the walls, while the wind howled through the trees, creating a symphony of nature's raw power that seemed both terrifying and mesmerizing at the same time.",
	"Walking through the cobblestone streets of the ancient town felt like stepping back in time, with every corner revealing a story etched into the walls of weathered buildings. The faint echoes of footsteps and the distant sound of a bell chiming added to the atmosphere, while the scent of freshly baked bread and blooming flowers from nearby markets made the journey feel like a sensory experience from a different era altogether.",
	"The smell of freshly baked bread wafted from the little bakery on the corner, mingling with the crisp morning air and making it impossible to resist stepping inside. The warm glow of the shop’s interior, filled with rows of pastries and loaves of bread, was inviting, and as the baker greeted her with a kind smile, she couldn’t help but feel as though she had discovered a slice of heaven tucked away in the heart of the bustling city.",
	"The vibrant colors of the sunset painted the sky in shades of orange, pink, and purple, creating a masterpiece that seemed too beautiful to be real. The reflection of the vivid hues danced on the surface of the calm lake, while the gentle chirping of crickets signaled the transition from day to night, wrapping the moment in an almost magical serenity that seemed to pause time itself",
	"He spent hours gazing at the star-studded night sky, marveling at the infinite universe and wondering how small and fleeting his own existence must seem in the grand scheme of things. With each shooting star that streaked across the heavens, he made silent wishes, dreaming of a future where he could explore the unknown and uncover the mysteries that had fascinated humankind since the dawn of time.",
}

// GetRandomText for generating the random coherent text for race
func GetRandomText() string {
	randomNum := rand.N(len(quotes))
	return quotes[randomNum]
}

type Stats struct {
	Accuracy float64
	// Word Per Minute
	WPM float64
}

func textDiffRatio(orig, newer []rune) float64 {
	if len(orig) == 0 || len(newer) == 0 {
		return 0
	}

	wrongs := 0
	for i := range orig {
		if i >= len(newer) {
			break
		}
		if orig[i] != newer[i] {
			wrongs++
		}
	}

	return 1 - float64(wrongs)/float64(len(orig))
}

// GetStats orig and user input text and total time in second
func GetStats(orig, newer []rune, sec float64) Stats {
	diffRatio := textDiffRatio(orig, newer)
	lpm := float64(len(newer)*60) / sec

	return Stats{
		Accuracy: diffRatio * 100,
		WPM:      lpm,
	}
}
