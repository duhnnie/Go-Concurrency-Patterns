package main

import (
	"fmt"
	"strings"
)

func generator(sentences []string) <-chan string {
	output := make(chan string)

	go func() {
		for _, sentence := range sentences {
			output <- sentence
		}

		close(output)
	}()

	return output
}

func tokenizer(input <-chan string, output chan<- string) {
	for sentence := range input {
		words := strings.Split(sentence, " ")

		for _, word := range words {
			output <- word
		}
	}

	close(output)
}

func filter(input <-chan string, output chan<- string) {
	for word := range input {
		if len(word) >= 4 {
			output <- word
		}
	}

	close(output)
}

func uppercaser(input <-chan string, output chan<- string) {
	for word := range input {
		output <- strings.ToUpper(word)
	}

	close(output)
}

func collector(input <-chan string, output chan<- []string) {
	words := []string{}

	for word := range input {
		words = append(words, word)
	}

	output <- words
	close(output)
}

func main() {
	var sentences = []string{
		"The cat jumped onto the table and knocked over a glass of milk",
		"She loves to read in the garden, surrounded by colorful flowers and birds singing all around her",
		"A small dog wagged its tail happily in the sunny park outside",
		"On weekends, he enjoys hiking in the mountains with his friends and exploring new trails together",
		"The sun set beautifully over the horizon, painting the sky in orange and pink hues",
		"They decided to bake cookies at home, filling the kitchen with sweet, delightful aromas and laughter",
		"A child flew a kite, watching it soar high in the bright blue sky",
		"The teacher explained the lesson clearly, helping students understand challenging topics and encouraging their questions and discussions",
		"The sound of rain on the roof is soothing during quiet evenings at home",
		"She took her camera to the beach and captured the stunning sunset with vibrant colors reflecting on the water",
		"A gentle breeze rustled the leaves of the tall trees in the park nearby",
		"At the library, they found many interesting books that made their hearts feel excited about learning and discovering new things",
		"He walked his bicycle along the path, enjoying the fresh air and bright sunshine",
		"The children played hide and seek in the backyard, giggling while trying to find good hiding spots from each other",
		"An artist painted a mural on the wall, bringing life and color to the once dull street",
		"During summer, families often go picnics at the lake, enjoying delicious food and playing games under the warm sun",
		"The puppy chased its tail in circles, making everyone laugh at its cute antics",
		"They visited the museum to learn about ancient history, seeing artifacts that sparked their curiosity and imagination",
		"Freshly baked bread filled the kitchen with a delicious aroma that made everyone hungry",
		"On rainy days, she enjoyed cozying up with her favorite blanket and reading adventure stories that took her to new places",
		"A butterfly landed softly on a flower, showcasing its beautiful colors in the sunlight",
		"The team practiced hard every day, hoping to win the championship and celebrate their success together with friends and families",
		"The little girl drew pictures of her family, expressing her love through colorful crayons and smiles",
		"In the morning, the coffee shop buzzed with chatter as people gathered to enjoy warm drinks and chat with friends",
		"The ice cream truck played a cheerful tune as it slowly drove down the street, attracting excited children everywhere it passed",
		"Many people enjoy gardening, taking pleasure in nurturing plants and watching them grow into beautiful blooms throughout the seasons",
		"The wind played softly with the wind chimes hanging in the porch, creating a melodic sound that echoed through the quiet evening",
		"He opened the window to let in fresh air, enjoying the feeling of nature surrounding him and filling the room",
		"The bright stars twinkled in the night sky, illuminating the dark world below with their tiny, shimmering lights",
		"They decided to volunteer at the local animal shelter, helping care for the animals in need and finding them loving homes.",
	}

	collectorChannel := make(chan []string)
	uppercaserChannel := make(chan string)
	filterChannel := make(chan string)
	tokenizerChannel := make(chan string)
	generatorChannel := generator(sentences)

	go collector(uppercaserChannel, collectorChannel)
	go uppercaser(filterChannel, uppercaserChannel)
	go filter(tokenizerChannel, filterChannel)
	go tokenizer(generatorChannel, tokenizerChannel)

	fmt.Println(<-collectorChannel)

}
