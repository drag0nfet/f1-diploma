export function initMenu() {
    const menuBtn = document.getElementById("menu-btn");
    const leftSidebar = document.getElementById("left-sidebar");
    const overlay = document.getElementById("overlay");

    menuBtn.addEventListener("click", function () {
        leftSidebar.classList.add("active");
        overlay.classList.add("active");
    });

    overlay.addEventListener("click", function () {
        leftSidebar.classList.remove("active");
        overlay.classList.remove("active");
    });
}