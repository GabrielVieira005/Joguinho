
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

// Inicia um buraco temporal com canais próprios
func iniciarBuraco(x, y int, jogo *Jogo) {
    state := &BuracoState{
        posX:    x,
        posY:    y,
        visivel: false,
        under:   Vazio, // começa invisível sobre chão vazio
    }

    // Goroutine de ciclo do buraco (aparece/desaparece periodicamente)
    go func() {
        for {
            // Fica invisível por 2-4 segundos
            tempoInvisivel := time.Duration(2+rand.Intn(3)) * time.Second
            time.Sleep(tempoInvisivel)
            
            // Aparece
            if !state.visivel {
                state.under = jogo.Mapa[state.posY][state.posX]
                jogo.Mapa[state.posY][state.posX] = BuracoVisivel
                state.visivel = true
                interfaceDesenharJogo(jogo)
            }
            
            // Fica visível por 3-5 segundos
            tempoVisivel := time.Duration(3+rand.Intn(3)) * time.Second
            time.Sleep(tempoVisivel)
            
            // Desaparece
            if state.visivel {
                jogo.Mapa[state.posY][state.posX] = state.under
                state.visivel = false
                interfaceDesenharJogo(jogo)
            }
        }
    }()

    // Goroutine verificadora de queda do jogador
    go func() {
        for {
            time.Sleep(100 * time.Millisecond)
            
            // Verifica se jogador está sobre o buraco quando ele está visível
            if state.visivel && jogo.PosX == state.posX && jogo.PosY == state.posY {
                canalQueda <- true
            }
        }
    }()
}
