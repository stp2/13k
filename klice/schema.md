# URL

|URL|popis|
|---|-----|
|/qr/... |url na QR kódy, povede na šifry|
|/signin| přihlášení|
|/team|info o týmu|
|/admin|admninistrace|
|/admin/teams|seznam týmů|
|/admin/reset|reset počítadla poslední šifry|
|/admin/tasks/'tier'|seznam šifer určené obtížnosti|

# Redis
|Klíč|Popis|
|----|-----|
|admin/user/'user'|heslo admin uživatele|
|start/'tier'|kde (QR) začíná jaká obtížnost|
|team/'team pass'/name|jméno|
|team/'team pass'/tier|obtížnost|
|team/'team pass'/last|poslední vyluštěná šifra|
|qr/'qr id'/position|poloha|
|qr/'qr id'/help|nápověda k poloze|
|qr/'qr id'/tier/'tier'|jakou šifru pro jakou obtížnost|
|qr/'qr id'/clue|nápověda do celkové soutěže|
|task/'task id'/number|číslo v pořadí|
|task/'task id'/cipher|šifra|
|task/'task id'/solution|řešení|
|task/'task id'/next|další QR kód|
|end/'end id'/position|poloha konce|
|end/'end id'/help|nápověda k poloze|
