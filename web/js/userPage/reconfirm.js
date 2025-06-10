import {showNotification} from "../notification";

export async function reconfirm() {
    const confirmationContent = document.getElementById("confirmation");
    confirmationContent.style.display = "block";

    const confirmationBtn = document.getElementById("send-confirmation-btn");

    confirmationBtn.addEventListener("click", async function (event) {
        event.preventDefault();
        const username = window.location.pathname.split('/').pop();
        const response = await fetch(`/reconfirmation?username=${username}`, {
            method: 'POST',
            credentials: 'include',
            headers: {
                'X-Requested-With': 'XMLHttpRequest'
            }
        });
        let data = await response.json();
        if (data.success) {
            showNotification(
                "success",
                "Проверьте электронную почту!"
            )
        }
        else {
            showNotification(
                "error",
                data.message
            )
        }
    });


}