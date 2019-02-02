# ToDo

* Si je fais un F5 sur playlist ou local musics, j'obtiens des pages vierges !
* Mettre une info pour FF pour le login pour bien autoriser avec un certif autosigné
* À la déconnexion, virer le jeton jwt
* Permettre de trier toutes les listes par ordre alpha
* Filtrer les affichages avec une partie recherche
* Regarder si je ne peux pas augmenter/décrémenter le son avec le molette de la souris
* Pour les ws:

```nginx
location /api/websocket {
proxy_pass http://portainer:9000;
proxy_http_version 1.1;
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection "Upgrade";
}
```