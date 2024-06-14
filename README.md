# LeTrèsBonCoin

## Description du projet

LeTrèsBonCoin est un réseau social de type forum 

## Nécessaire

- Docker
- Golang
- Docker Compose


## Utilisation du projet

Importez le répertoire localement :

```bash
git clone https://github.com/Wawanke/LeTresBonCoin.git 
```

##  Lancer les dockers

Installation:
git pull https://github.com/Wawanke/LeTresBonCoin.git
ensuit une fois dedans 

Voici les étapes à suivre après avoir récupérer le projet :

Il faut instancier le dock contenant le back en Golang :

```bash
docker build -t golang .
```

Ensuite la même avec le docker contenant le front en HTML/CSS :

```bash
docker build -t html .
```

Après la même avec le docker contenant la base de donnée :

En allant dans le dossier DbDocker ouvrir un bash :

```bash
docker build -t  db .
```

## Kervoelen Erwann, Joël Rakotomalala