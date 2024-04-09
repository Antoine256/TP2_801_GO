package global

import (
	. "github.com/pspaces/gospace/space"
)

func Add[T any](ts *Space, name string, value T) {
	var x T
	ts.Put(name, value) // On initialise le tuple si jamais il n'existe pas
	ts.GetAll(name, &x) // On récupère tous les tuples pour être sûr de ne pas avoir de doublons
	ts.Put(name, value) // On met à jour la nouvelle valeur
}
