package main

import (
	"math/rand"
	"time"
)

type BuracoState struct {
	posX, posY int
	visivel    bool
	under      Elemento
}

var canalQueda = make(chan bool)

// Inicia um buraco temporal com canais próprios e timeout para queda do jogador
func iniciarBuraco(x, y int, jogo *Jogo) {
	state := &BuracoState{
		posX:    x,
		posY:    y,
		visivel: false,
		under:   Vazio, // começa invisível sobre chão vazio
	}

	go func() {
		for {
			// Fica invisível por 2-4 segundos
			tempoInvisivel := time.Duration(2+rand.Intn(3)) * time.Second
			time.Sleep(tempoInvisivel)

			// Aparece
			<-jogoMutex
			state.under = jogo.Mapa[state.posY][state.posX]
			jogo.Mapa[state.posY][state.posX] = BuracoVisivel
			state.visivel = true
			interfaceDesenharJogo(jogo)
			jogoMutex <- struct{}{}

			caiuChan := make(chan struct{})
			done := make(chan struct{})

			// Goroutine para monitorar queda enquanto o buraco está visível
			go func() {
				defer close(done)
				for state.visivel {
					time.Sleep(100 * time.Millisecond)
					if jogo.PosX == state.posX && jogo.PosY == state.posY {
						select {
						case caiuChan <- struct{}{}:
						default:
						}
						return
					}
				}
			}()

			select {
			case <-caiuChan:
				<-jogoMutex
				canalQueda <- true
				state.visivel = false
				jogo.Mapa[state.posY][state.posX] = state.under
				interfaceDesenharJogo(jogo)
				jogoMutex <- struct{}{}
				time.Sleep(1 * time.Second)
			case <-time.After(3 * time.Second):
				<-jogoMutex
				state.visivel = false
				jogo.Mapa[state.posY][state.posX] = state.under
				jogo.StatusMsg = "O buraco sumiu sozinho!"
				interfaceDesenharJogo(jogo)
				jogoMutex <- struct{}{}
				time.Sleep(1 * time.Second)
				<-jogoMutex
				jogo.StatusMsg = ""
				jogoMutex <- struct{}{}
			}
			<-done // Aguarda a goroutine terminar antes de seguir para o próximo ciclo
		}
	}()
}
