package main

import (
	"time"
)

var armadilhaAtiva = make(map[[2]int]bool)

// Inicia uma armadilha temporizada com canal e timeout
func iniciarArmadilha(x, y int, jogo *Jogo) {
	pos := [2]int{x, y}
	if armadilhaAtiva[pos] {
		return // já existe armadilha ativa nesta posição
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

			select {
			case <-ativou:
				canalArmadilha <- true
				jogo.PresoAte = time.Now().Add(2 * time.Second) // prende por 2 segundos
				jogo.Mapa[y][x] = Vazio
				interfaceDesenharJogo(jogo)
				time.Sleep(1 * time.Second)
			case <-time.After(3 * time.Second):
				jogo.Mapa[y][x] = Vazio
				jogo.StatusMsg = "A armadilha desarmou sozinha!"
				interfaceDesenharJogo(jogo)
				time.Sleep(1 * time.Second)
				jogo.StatusMsg = ""
			}
			<-done
			// Aguarda 3 segundos antes de reativar a armadilha
			time.Sleep(3 * time.Second)
			// Só reativa se o local ainda estiver vazio (não foi sobrescrito por outro elemento)
			if jogo.Mapa[y][x].simbolo != Vazio.simbolo {
				break
			}
		}
		armadilhaAtiva[pos] = false
	}()
}
