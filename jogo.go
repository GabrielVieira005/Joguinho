// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"math/rand"
	"os"
	"time"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo  rune
	cor      Cor
	corFundo Cor
	tangivel bool // Indica se o elemento bloqueia passagem
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa           [][]Elemento // grade 2D representando o mapa
	PosX, PosY     int          // posição atual do personagem
	UltimoVisitado Elemento     // elemento que estava na posição do personagem antes de mover
	StatusMsg      string       // mensagem para a barra de status
	PresoAte       time.Time    // se não zero, personagem está preso até esse instante
}

// Elementos visuais do jogo
var (
	Personagem    = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo       = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede        = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao     = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio         = Elemento{' ', CorPadrao, CorPadrao, false}
	Patrulheiro   = Elemento{'⚉', CorAmarelo, CorPadrao, true}
	BuracoVisivel = Elemento{'●', CorVermelho, CorPadrao, false}
	Armadilha     = Elemento{'^', CorMagenta, CorPadrao, false}
)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo{UltimoVisitado: Vazio}
}

// Canal para armadilha temporizada
var canalArmadilha = make(chan bool)

// Canal para exclusão mútua (mutex por canal)
var jogoMutex = make(chan struct{}, 1)

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case Inimigo.simbolo:
				e = Inimigo
			case Vegetacao.simbolo:
				e = Vegetacao
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y // registra a posição inicial do personagem
			case Patrulheiro.simbolo: //Modificação: iniciar patrulheiro
				e = Vazio
				iniciarPatrulheiroAsync(x, y, jogo)
			case '●': // ADICIONE ESTAS LINHAS - Buraco temporal
				e = Vazio
				iniciarBuraco(x, y, jogo)
			case '^': // Armadilha temporizada
				e = Vazio // <-- Corrija aqui: nunca desenhe ^ diretamente
				iniciarArmadilha(x, y, jogo)
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	<-jogoMutex
	defer func() { jogoMutex <- struct{}{} }()
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	<-jogoMutex
	defer func() { jogoMutex <- struct{}{} }()
	nx, ny := x+dx, y+dy

	// Obtem elemento atual na posição
	elemento := jogo.Mapa[y][x] // guarda o conteúdo atual da posição

	jogo.Mapa[y][x] = jogo.UltimoVisitado   // restaura o conteúdo anterior
	jogo.UltimoVisitado = jogo.Mapa[ny][nx] // guarda o conteúdo atual da nova posição
	jogo.Mapa[ny][nx] = elemento            // move o elemento
}

// Inicia a goroutine do patrulheiro em uma posição aleatória
func iniciarPatrulheiro(x, y int, jogo *Jogo) {
	go func() {
		for {
			dx, dy := 0, 0
			// Move aleatoriamente: -1, 0 ou 1
			switch rand.Intn(4) {
			case 0:
				dx = -1
			case 1:
				dx = 1
			case 2:
				dy = -1
			case 3:
				dy = 1
			}

			// Verifica nova posição
			nx, ny := x+dx, y+dy
			<-jogoMutex
			if jogoPodeMoverPara(jogo, nx, ny) {
				jogoMoverElemento(jogo, x, y, dx, dy)
				x, y = nx, ny
			}
			jogoMutex <- struct{}{}
			time.Sleep(500 * time.Millisecond) // Aguarda meio segundo
		}
	}()
}

// Inicia a goroutine do buraco temporal
// (implementação removida para evitar conflito de declaração; veja buraco.go)

// Inicia a goroutine do buraco temporal
// (implementação removida para evitar conflito de declaração; veja buraco.go)

// Mapa global para canais dos portais ativos
var portalChannels = make(map[[2]int]chan bool)

// Registra canal do portal na posição (x, y)
func registerPortalChannel(x, y int, ch chan bool) {
	portalChannels[[2]int{x, y}] = ch
}

// Remove canal do portal na posição (x, y)
func unregisterPortalChannel(x, y int) {
	delete(portalChannels, [2]int{x, y})
}
