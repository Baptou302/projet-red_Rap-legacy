🎤 Rap-Legacy

Rap-Legacy est un jeu vidéo codé en Go (Golang) et développé en une semaine.
Il plonge le joueur dans l’univers du rap game, retraçant l’histoire d’un jeune rappeur en quête de reconnaissance dans l’industrie musicale.

🚀 Présentation du projet

Le jeu reprend des mécaniques de combat au tour par tour, où chaque attaque correspond à une punchline et où l’égo du rappeur représente sa barre de vie.
À travers combats, progression et gestion de ressources, le joueur évolue jusqu’à devenir une véritable légende du micro.

🎮 Fonctionnalités

- Combat au tour par tour avec système d’égo (barre de vie).

- Classes de rappeurs sélectionnables en début de partie.

- Système de followers servant de points d’expérience.

- Inventaire avec items consommables offrant des bonus.

- Crafting pour combiner des objets et en créer de nouveaux.

- Économie et marchand pour acheter et gérer ses ressources.

- Sauvegarde fonctionnelle à chaque lancement du jeu.

- Bande-son originale pour renforcer l’immersion.

👥 Répartition des tâches

Même si toutes les fonctionnalités ont été conçues, testées et finalisées en duo, voici la répartition des principales contributions :

Baptiste

- Réalisation des images et assets visuels

- Conception et intégration des menus

- Développement du système de combats

- Intégration de la musique

Erwan : 

- Système de sauvegarde

- Développement de l’inventaire

- Mise en place du marchand

- Implémentation du système de craft

Ensemble :

- Conception et implémentation des potions

- Travail collaboratif sur l’ensemble du code et des fonctionnalités

⚙️ Défis rencontrés

- Temps limité : une semaine seulement pour réaliser le jeu.

- Problèmes de merge Git : nécessitant la mise en place d’un protocole interne pour limiter les conflits.

- Architecture du projet : relier tous les fichiers entre eux et organiser une arborescence cohérente.

🛠️ Technologies utilisées

Langage : Golang

Gestion de version : Git & GitHub / Ebiten

📦 Installation et lancement

Clonez le dépôt :

git clone https://github.com/Baptou302/projet-red_Rap-legacy.git

Une fois cloner entrez cette commande : 

cd rap-legacy

 Pour Lancez le jeu entrez cette commande :

go run main.go / ou alors / go run .