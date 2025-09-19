// main.go - Loop principal do jogo
package main

import (
	"os"
	"time"
	"math/rand"//Usado para teletransporte (buraco temporal)
)

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

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

	// Modificação: Goroutine para detectar morte por patrulheiro
    go func() {
		for range canalMorte {
			jogo.StatusMsg = "Você foi pego pelo patrulheiro!"
			interfaceDesenharJogo(&jogo)
			time.Sleep(2 * time.Second)
			os.Exit(0)
		}
    }()

	// Modificação: Goroutine para detectar queda no buraco
	go func() {
        for range canalQueda {
            jogo.StatusMsg = "Você caiu em um buraco temporal!"
            interfaceDesenharJogo(&jogo)
            time.Sleep(1 * time.Second)
            
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
            
            // Limpa a mensagem após um tempo
            time.Sleep(2 * time.Second)
            jogo.StatusMsg = ""
        }
    }()

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}
		
		interfaceDesenharJogo(&jogo)
	}
}