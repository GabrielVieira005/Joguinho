// main.go - Loop principal do jogo
package main

import (
	"math/rand" //Usado para teletransporte (buraco temporal)
	"os"
	"time"
)

func main() {
	// Inicializa o mutex
	jogoMutex <- struct{}{}

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Loop principal de entrada
	for {
		select {
		case <-canalMorte:
			<-jogoMutex
			jogo.StatusMsg = "Você foi pego pelo patrulheiro!"
			interfaceDesenharJogo(&jogo)
			jogoMutex <- struct{}{}
			time.Sleep(2 * time.Second)
			os.Exit(0)
		case <-canalQueda:
			<-jogoMutex
			jogo.StatusMsg = "Você caiu em um buraco temporal!"
			interfaceDesenharJogo(&jogo)

			// Teletransporta jogador para posição aleatória válida
			for tentativas := 0; tentativas < 100; tentativas++ {
				newX := rand.Intn(len(jogo.Mapa[0]))
				newY := rand.Intn(len(jogo.Mapa))

				if jogoPodeMoverPara(&jogo, newX, newY) {
					jogo.PosX, jogo.PosY = newX, newY
					jogo.StatusMsg = "Você foi teletransportado para outro local!"
					break
				}
			}
			interfaceDesenharJogo(&jogo)
			jogoMutex <- struct{}{}
			time.Sleep(2 * time.Second)
			<-jogoMutex
			jogo.StatusMsg = ""
			jogoMutex <- struct{}{}
		case <-canalArmadilha:
			<-jogoMutex
			jogo.StatusMsg = "Você caiu em uma armadilha!"
			interfaceDesenharJogo(&jogo)
			jogoMutex <- struct{}{}
			time.Sleep(2 * time.Second)
			<-jogoMutex
			jogo.StatusMsg = ""
			jogoMutex <- struct{}{}
		default:
			evento := interfaceLerEventoTeclado()
			<-jogoMutex
			if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
				jogoMutex <- struct{}{}
				break
			}
			interfaceDesenharJogo(&jogo)
			jogoMutex <- struct{}{}
			time.Sleep(10 * time.Millisecond)
		}
	}
}
}
