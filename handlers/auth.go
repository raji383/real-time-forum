package handlers

import (
    "html/template"
    "net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("template/index.html")
    if err != nil {
        http.Error(w, "Template error", http.StatusInternalServerError)
        return
    }

    // يمكن لاحقًا تمرير بيانات المستخدم أو المنشورات هنا
    tmpl.Execute(w, nil)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/", http.StatusSeeOther) // نفس الصفحة فيها login
}

func SignupPage(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Login(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
   // username := r.FormValue("username")
   // password := r.FormValue("password")

    // TODO: تحقق من المستخدم في قاعدة البيانات

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Signup(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
  //  username := r.FormValue("username")
  //  email := r.FormValue("email")
    pass1 := r.FormValue("pass1")
    pass2 := r.FormValue("pass2")

    if pass1 != pass2 {
        http.Error(w, "Passwords do not match", http.StatusBadRequest)
        return
    }

    // TODO: إدخال المستخدم لقاعدة البيانات

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
    // TODO: حذف الكوكي أو التوكن
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
