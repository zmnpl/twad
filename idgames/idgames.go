package idgames

const (
	API_URL = "https://www.doomworld.com/idgames/api/api.php?"
	API_GET = "action=get&id=2319&out=json"
)

type Idgame struct {
	Id          int      `json:"id"`          // The file's id.
	Title       string   `json:"title"`       // The title of the file.
	Dir         string   `json:"dir"`         // The file's full directory path.
	Filename    string   `json:"filename"`    // The filename itself, no path.
	Size        int      `json:"size"`        // The size of the file in bytes.
	Age         int64    `json:"age"`         // The date that the file was added in seconds since the Unix Epoch (Jan. 1, 1970). Note: This is likely influenced by the time zone of the primary idGames Archive.
	Date        string   `json:"date"`        // A YYYY-MM-DD formatted date describing the date that this file was added to the archive.
	Author      string   `json:"author"`      // The file's author/uploader.
	Email       string   `json:"email"`       // The author's E-mail address.
	Description string   `json:"description"` // The file's description.
	Credits     string   `json:"credits"`     // The file's additional credits.
	Base        string   `json:"base"`        // The file's base (from another mod? made from scratch?).
	Buildtime   string   `json:"buildtime"`   // The file's/WAD's build time.
	Editors     string   `json:"editors"`     // The editors used to create this.
	Bugs        string   `json:"bugs"`        // Known bugs (if any).
	Textfile    string   `json:"textfile"`    // The file's text file contents.
	Rating      int      `json:"rating"`      // The file's average rating, as rated by users.
	Votes       int      `json:"votes"`       // The number of votes that this file received.
	Url         string   `json:"url"`         // The URL for the idGames Archive page for this file.
	Idgamesurl  string   `json:"idgamesurl"`  // The idgames protocol URL for this file.
	Reviews     []Review `json:"reviews"`     // The element that contains all reviews for this file in review elements.
}

type Review struct {
	Text     string `json:"text"`     // The individual review's text, if any. Note: may be blank.
	Vote     int    `json:"vote"`     // The vote associated with the review.
	Username string `json:"username"` // The user name associated with the review, if any. Note: may be blank/null, which means "Anonymous". Since Version 3
}
