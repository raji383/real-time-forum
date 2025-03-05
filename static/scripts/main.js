document.addEventListener("DOMContentLoaded", function () {


    document.querySelectorAll(".post").forEach(post => {
        let likeBtn = post.querySelector(".like");
        let dislikeBtn = post.querySelector(".dislike");
        let likeSpan = likeBtn.querySelector("span");
        let dislikeSpan = dislikeBtn.querySelector("span");
        let postId = post.getAttribute("data-post-id");

        likeBtn.addEventListener("click", async function () {
            let likes = Number(likeSpan.textContent);
            let dislikes = Number(dislikeSpan.textContent);

            if (likeBtn.classList.contains("active")) {
                let success = await updateLikes(postId, "like", -1);
                if (success) {
                    likeBtn.classList.remove("active");
                    likeSpan.textContent = likes - 1;
                }
            } else {
                let success = await updateLikes(postId, "like", 1);
                if (success) {
                    likeBtn.classList.add("active");
                    likeSpan.textContent = likes + 1;

                    if (dislikeBtn.classList.contains("active")) {
                        await updateLikes(postId, "dislike", -1);
                        dislikeBtn.classList.remove("active");
                        dislikeSpan.textContent = dislikes - 1;
                    }
                }
            }
        });

        dislikeBtn.addEventListener("click", async function () {
            let likes = Number(likeSpan.textContent);
            let dislikes = Number(dislikeSpan.textContent);

            if (dislikeBtn.classList.contains("active")) {
                let success = await updateLikes(postId, "dislike", -1);
                if (success) {
                    dislikeBtn.classList.remove("active");
                    dislikeSpan.textContent = dislikes - 1;
                }
            } else {
                let success = await updateLikes(postId, "dislike", 1);
                if (success) {
                    dislikeBtn.classList.add("active");
                    dislikeSpan.textContent = dislikes + 1;

                    if (likeBtn.classList.contains("active")) {
                        await updateLikes(postId, "like", -1);
                        likeBtn.classList.remove("active");
                        likeSpan.textContent = likes - 1;
                    }
                }
            }
        });
    })

    document.querySelectorAll(".main-comment").forEach(but => {
        let button = but.querySelector("#showComment")
        let n = 0
        button.addEventListener("click", function () {
            if (n === 0) {
                but.querySelector(".comment-text").style.display = "flex"
                button.innerHTML = '<i class="fa-solid fa-comment-slash"></i>'
                n = 1
            } else {
                but.querySelector(".comment-text").style.display = "none"
                button.innerHTML = '<i class="fa-solid fa-comment"></i>'
                n = 0
            }
        })

        let commentText = but.querySelector("#com")
        let text = but.querySelector("input")
        let username = document.getElementById("username")
        let u = 0
        if (username === null) {
            u = 1
        } else {
            username = username.textContent
        }
        let postID = but.getAttribute("data-post-id")


        commentText.addEventListener("click", function () {
            if (text.value === "") {
                return
            }
            document.querySelectorAll(".comment-content").forEach(async function (el) {



                let nameOfUser = document.createElement("h3")
                let content = document.createElement("pre")
                let time = document.createElement("span")
                let div = document.createElement("div")
                if (u === 1) {
                    location.href = "/login-page"
                    return
                }
                let date = new Date()
                let arrDate = date.toISOString().split("T")
                let dateNow = arrDate[0] + " " + arrDate[1].slice(0, 8)


                nameOfUser.textContent = username
                content.textContent = text.value
                time.textContent = dateNow

                div.className = "div"
                div.append(nameOfUser, content, time)


                let elID = el.getAttribute("data-post-id")
                if (elID === postID) {
                    text.value = ""
                    await addComment(postID, nameOfUser.textContent, content.textContent, dateNow)
                    el.querySelector(".hihi").prepend(div)
                }

            })

        })

    })

    document.querySelectorAll(".comment-content").forEach(el => {
        let c = 0
        el.querySelector(".hide").addEventListener("click", function () {
            if (c === 0) {
                el.querySelector(".hihi").style.display = "block"
                this.textContent = "Hide Comments"
                c = 1
            } else {
                el.querySelector(".hihi").style.display = "none"
                this.textContent = "Show Comments"
                c = 0
            }
        })
    })


});

async function addComment(postID, username, content, time) {
    try {
        let response = await fetch("/comment", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ post_id: Number(postID), nameOfUser: username, comment: content, time: time })
        })

        let data = await response.json();
        console.log(data);

        return true
    } catch (error) {
        console.error(error);
        return false
    }

}

async function updateLikes(postId, action, change) {
    try {
        let response = await fetch("/like", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ post_id: Number(postId), action: action, change: change })
        });

        let data = await response.json();

        if (!data.success) {
            window.location.href = "/login-page"
            return false
        }

        return true

    } catch (error) {
        console.error("Fetch error:", error);
        return false
    }
}


function openNav() {
    document.getElementById("mySidenav").style.width = "250px";
}
function closeNav() {
    document.getElementById("mySidenav").style.width = "0";
}

// let c = 0

// function showComments() {
//     if (c === 0) {
//         document.querySelector(".comment-content").style.display = "block"
//         document.querySelector(".hide").textContent = "Hide comments"
//         c = 1
//     } else {
//         document.querySelector(".comment-content").style.display = "none"
//         document.querySelector(".hide").textContent = "Show comments"
//         c = 0
//     }

// }