// Populate the sidebar
//
// This is a script, and not included directly in the page, to control the total size of the book.
// The TOC contains an entry for each page, so if each page includes a copy of the TOC,
// the total size of the page becomes O(n**2).
class MDBookSidebarScrollbox extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        this.innerHTML = '<ol class="chapter"><li class="chapter-item "><a href="nova.html"><strong aria-hidden="true">1.</strong> Nova</a></li><li class="chapter-item "><a href="quickstart.html"><strong aria-hidden="true">2.</strong> Quickstart</a></li><li class="chapter-item "><div><strong aria-hidden="true">3.</strong> Features</div><a class="toggle"><div>❱</div></a></li><li><ol class="section"><li class="chapter-item "><a href="cli.html"><strong aria-hidden="true">3.1.</strong> CLI</a></li><li class="chapter-item "><a href="scaffolding.html"><strong aria-hidden="true">3.2.</strong> Scaffolding</a></li><li class="chapter-item "><a href="router.html"><strong aria-hidden="true">3.3.</strong> Router</a></li><li class="chapter-item "><a href="database_migrations.html"><strong aria-hidden="true">3.4.</strong> Database Migrations</a></li><li class="chapter-item "><a href="openapi.html"><strong aria-hidden="true">3.5.</strong> OpenAPI Spec</a></li></ol></li><li class="chapter-item "><a href="middleware.html"><strong aria-hidden="true">4.</strong> Middleware</a></li><li class="chapter-item "><a href="design.html"><strong aria-hidden="true">5.</strong> Design</a></li><li class="chapter-item "><a href="faq.html"><strong aria-hidden="true">6.</strong> FAQ</a></li><li class="chapter-item "><a href="reference.html"><strong aria-hidden="true">7.</strong> Reference</a></li></ol>';
        // Set the current, active page, and reveal it if it's hidden
        let current_page = document.location.href.toString().split("#")[0].split("?")[0];
        if (current_page.endsWith("/")) {
            current_page += "index.html";
        }
        var links = Array.prototype.slice.call(this.querySelectorAll("a"));
        var l = links.length;
        for (var i = 0; i < l; ++i) {
            var link = links[i];
            var href = link.getAttribute("href");
            if (href && !href.startsWith("#") && !/^(?:[a-z+]+:)?\/\//.test(href)) {
                link.href = path_to_root + href;
            }
            // The "index" page is supposed to alias the first chapter in the book.
            if (link.href === current_page || (i === 0 && path_to_root === "" && current_page.endsWith("/index.html"))) {
                link.classList.add("active");
                var parent = link.parentElement;
                if (parent && parent.classList.contains("chapter-item")) {
                    parent.classList.add("expanded");
                }
                while (parent) {
                    if (parent.tagName === "LI" && parent.previousElementSibling) {
                        if (parent.previousElementSibling.classList.contains("chapter-item")) {
                            parent.previousElementSibling.classList.add("expanded");
                        }
                    }
                    parent = parent.parentElement;
                }
            }
        }
        // Track and set sidebar scroll position
        this.addEventListener('click', function(e) {
            if (e.target.tagName === 'A') {
                sessionStorage.setItem('sidebar-scroll', this.scrollTop);
            }
        }, { passive: true });
        var sidebarScrollTop = sessionStorage.getItem('sidebar-scroll');
        sessionStorage.removeItem('sidebar-scroll');
        if (sidebarScrollTop) {
            // preserve sidebar scroll position when navigating via links within sidebar
            this.scrollTop = sidebarScrollTop;
        } else {
            // scroll sidebar to current active section when navigating via "next/previous chapter" buttons
            var activeSection = document.querySelector('#sidebar .active');
            if (activeSection) {
                activeSection.scrollIntoView({ block: 'center' });
            }
        }
        // Toggle buttons
        var sidebarAnchorToggles = document.querySelectorAll('#sidebar a.toggle');
        function toggleSection(ev) {
            ev.currentTarget.parentElement.classList.toggle('expanded');
        }
        Array.from(sidebarAnchorToggles).forEach(function (el) {
            el.addEventListener('click', toggleSection);
        });
    }
}
window.customElements.define("mdbook-sidebar-scrollbox", MDBookSidebarScrollbox);
