package main

import (
	"time"
)
// Canais da armadilha, tamanho e se o personagem pisou o não
var armadilhaAtiva = make(map[[2]int]bool)
var canalArmadilha = make(chan bool)

// Inicia uma armadilha temporizada com canal e timeout
func iniciarArmadilha(x, y int, jogo *Jogo) {
	pos := [2]int{x, y}
	if armadilhaAtiva[pos] {
		return // Se já existe armadilha ativa nesta posição
	}
	armadilhaAtiva[pos] = true

	go func() {
		for {
			// Ativa a armadilha
			jogo.Mapa[y][x] = Armadilha
			interfaceDesenharJogo(jogo)

			// Canal para detectar ativação
			ativou := make(chan struct{})
			done := make(chan struct{})

			// Goroutine para monitorar se o jogador pisa na armadilha
			go func() {
				defer close(done)
				for jogo.Mapa[y][x].simbolo == Armadilha.simbolo {
					time.Sleep(100 * time.Millisecond)
					if jogo.PosX == x && jogo.PosY == y {
						select {
						case ativou <- struct{}{}:
						default:
						}
						return
					}
				}
			}()
			
			//Sistema que faz o o jogador ficar preso na armadilha
			select {
			case <-ativou:
				canalArmadilha <- true
				jogo.PresoAte = time.Now().Add(3 * time.Second) // prende por 3 segundos
				jogo.Mapa[y][x] = ArmadilhaUsada
				interfaceDesenharJogo(jogo)
				time.Sleep(7 * time.Second)
			case <-time.After(5 * time.Second):
				jogo.Mapa[y][x] = Vazio
				interfaceDesenharJogo(jogo)
				time.Sleep(1 * time.Second)
			}
			<-done
			
			// Aguarda 4 segundos antes de reativar a armadilha
			time.Sleep(4 * time.Second)
			
			// Só reativa se o local ainda estiver vazio ou com armadilha usada
			elemento := jogo.Mapa[y][x]
			if elemento.simbolo != Vazio.simbolo && elemento.simbolo != ArmadilhaUsada.simbolo {
				break
			}
		}
		armadilhaAtiva[pos] = false
	}()
}