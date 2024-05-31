package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

type Card struct {
	// Suit can be "Spades", "Hearts", "Diamonds", "Clubs"
	Suit string
	// Number can be 1-13, where 1 is ace, 11 is jack, 12 is queen, 13 is king
	Number int
}

func (c Card) String() string {
	var value string
	switch c.Number {
	case 1:
		value = "Ace"
	case 11:
		value = "Jack"
	case 12:
		value = "Queen"
	case 13:
		value = "King"
	default:
		value = fmt.Sprintf("%d", c.Number)
	}
	return fmt.Sprintf("%s of %s", value, c.Suit)
}

func printHand(hand []Card) string {
	var strs []string
	for _, card := range hand {
		strs = append(strs, card.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
}

// Value returns the value of the card in blackjack.
// For aces, always return the high value of 11.
// When calculating a player hand value, the aces will
// be adjusted to 1 if the hand value is over 21.
func (c Card) Value() int {
	if c.Number == 1 {
		return 11
	} else if c.Number >= 10 {
		return 10
	}
	return c.Number
}

type Deck struct {
	Cards []Card
}

func NewDeck() *Deck {
	deck := &Deck{}
	for _, suit := range []string{"♠", "♥", "♦", "♣"} {
		for i := 1; i <= 13; i++ {
			// suit = Hearts
			// i = 1
			deck.Cards = append(deck.Cards, Card{Suit: suit, Number: i})
		}
	}
	deck.Shuffle()
	return deck
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

func (d *Deck) Draw() Card {
	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return card
}

type Player struct {
	Name    string
	Hand    []Card
	Balance int
	Bet     int

	GamesPlayed int
	Wins        int
	Losses      int
	WinStreak   int
	LoseStreak  int

	PlayAlgorithm func(human *Player) string
	BetAlgorithm  func(human *Player) int
}

func NewPlayer(name string, startingBalance int, playAlgorithm func(*Player) string,
	betAlgorithm func(*Player) int) *Player {
	return &Player{
		Name:          name,
		Balance:       startingBalance,
		PlayAlgorithm: playAlgorithm,
		BetAlgorithm:  betAlgorithm,
	}
}

func (p *Player) HandValue() int {
	var value int
	var aceCount int
	for _, card := range p.Hand {
		// TODO: Always assuming aces are high
		value += card.Value()
		if card.Value() == 11 {
			aceCount++
		}
	}

	if value <= 21 {
		return value
	}

	for {
		if aceCount == 0 {
			break
		}

		value -= 10
		aceCount--
		if value <= 21 {
			break
		}
	}
	return value
}

func (p *Player) PlayHand(deck *Deck) {
	for {
		action := p.PlayAlgorithm(p)

		if action == "hit" {
			p.Hand = append(p.Hand, deck.Draw())
			fmt.Println(printHand(p.Hand))
			fmt.Println(p.HandValue())

			if p.HandValue() > 21 {
				fmt.Println("Bust!")
				return
			}
		} else if action == "stand" {
			return
		}
	}
}

func (p *Player) CompareWithDealer(dealer *Player) string {

	// If player busts, player loses
	if p.HandValue() > 21 {
		return "lose"
	}

	// If player doesn't bust and dealer busts, player wins
	if dealer.HandValue() > 21 {
		return "win"
	}

	// If player has higher value than the dealer, player wins
	if p.HandValue() > dealer.HandValue() {
		return "win"
	}

	// If player has equal value than the dealer, push
	if p.HandValue() == dealer.HandValue() {
		return "push"
	}

	// If player has lower value than the dealer, player loses
	if p.HandValue() < dealer.HandValue() {
		return "lose"
	}

	return "shouldn't happen"
}

func (p *Player) Reconcile(dealer *Player) {
	finalStatus := p.CompareWithDealer(dealer)

	p.GamesPlayed++

	if finalStatus == "win" {
		if p.HandValue() == 21 && len(p.Hand) == 2 {
			p.Balance += int(float64(p.Bet) * 2.5)
		} else {
			p.Balance += p.Bet * 2
		}
		p.Bet = 0
		p.WinStreak++
		p.Wins++
		p.LoseStreak = 0
	} else if finalStatus == "push" {
		p.Balance += p.Bet
		p.Bet = 0
	} else {
		p.Bet = 0
		p.LoseStreak++
		p.Losses++
		p.WinStreak = 0
	}
}

func (p *Player) PrintStatistics() string {
	return fmt.Sprintf("%s has played %d games, won %d games, and lost %d games, with a win streak of %d, and a lose streak of %d, balance: $%d",
		p.Name, p.GamesPlayed, p.Wins, p.Losses, p.WinStreak, p.LoseStreak, p.Balance)
}

func dealerPlayAlgorithm(dealer *Player) string {
	if dealer.HandValue() < 17 {
		return "hit"
	}
	return "stand"
}

func dealerBetAlgorithm(dealer *Player) int {
	return 0
}

func humanPlayAlgorithm(player *Player) string {
	fmt.Println("Would you like to 'hit' or 'stand'?")
	var input string
	fmt.Scanln(&input)
	return input
}

func humanBetAlgorithm(player *Player) int {
	fmt.Printf("%s, how much would you like to bet?\n", player.Name)
	var bet int
	fmt.Scanln(&bet)
	return bet
}

func brucePlayAlgorithm(player *Player) string {
	if player.HandValue() < 15 {
		return "hit"
	}
	return "stand"
}

func bruceBetAlgorithm(player *Player) int {
	var bet int
	if player.LoseStreak < 3 {
		bet = 5 * int(math.Pow(2, float64(player.LoseStreak)))
	} else {
		bet = 5
		player.LoseStreak = 0
	}
	return bet
}

func main() {
	fmt.Println("Welcome to Blackjack!")

	dealer := NewPlayer("Dealer", 0, dealerPlayAlgorithm, dealerBetAlgorithm)
	bruce := NewPlayer("Bruce", 100, brucePlayAlgorithm, bruceBetAlgorithm)
	human := NewPlayer("Human", 100, humanPlayAlgorithm, humanBetAlgorithm)

	for i := 0; i < 10000; i++ {
		deck := NewDeck()

		// Place bets
		bet := bruce.BetAlgorithm(bruce)
		bruce.Bet = bet
		bruce.Balance = bruce.Balance - bet

		bet = human.BetAlgorithm(human)
		human.Bet = bet
		human.Balance = human.Balance - bet

		// deal the cards
		for i := 0; i < 2; i++ {
			bruce.Hand = append(bruce.Hand, deck.Draw())
			human.Hand = append(human.Hand, deck.Draw())
			dealer.Hand = append(dealer.Hand, deck.Draw())
		}

		// Reveal dealer's hand
		fmt.Printf("\nDealer's hand: [%s] [x]\n", dealer.Hand[0].String())

		// First player
		fmt.Printf("\nIt's %s's turn\n", bruce.Name)
		fmt.Println(printHand(bruce.Hand))
		fmt.Println(bruce.HandValue())
		bruce.PlayHand(deck)

		// Second player
		fmt.Printf("\nIt's %s's turn\n", human.Name)
		fmt.Println(printHand(human.Hand))
		fmt.Println(human.HandValue())
		human.PlayHand(deck)

		// Dealer last
		fmt.Printf("\nIt's %s's turn\n", dealer.Name)
		fmt.Println(printHand(dealer.Hand))
		fmt.Println(dealer.HandValue())
		dealer.PlayHand(deck)

		fmt.Printf("Game #%d over!\n\n", i+1)

		// Reconcile the bets and win/loss counts
		bruce.Reconcile(dealer)
		human.Reconcile(dealer)

		fmt.Printf("%s %s: %d %s\n", bruce.Name, bruce.CompareWithDealer(dealer), bruce.HandValue(), printHand(bruce.Hand))
		fmt.Printf("%s %s: %d %s\n", human.Name, human.CompareWithDealer(dealer), human.HandValue(), printHand(human.Hand))

		fmt.Println()
		fmt.Println(bruce.PrintStatistics())
		fmt.Println(human.PrintStatistics())

		// Done playing?
		fmt.Println("Would you like to play again? 'yes' or 'no'")
		var input string
		fmt.Scanln(&input)
		if input == "no" {
			break
		}

		// Reset the hands
		bruce.Hand = []Card{}
		human.Hand = []Card{}
		dealer.Hand = []Card{}
	}
}
