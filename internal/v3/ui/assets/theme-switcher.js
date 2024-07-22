/*!
 * Minimal theme switcher
 *
 * Pico.css - https://picocss.com
 * Copyright 2019-2024 - Licensed under MIT
 */

const themeSwitcher = {
  // Config
  _scheme: "light",
  rootAttribute: "data-theme",
  localStorageKey: "picoPreferredColorScheme",

  // Init
  init() {
    this.scheme = this.schemeFromLocalStorage;
    this.initSwitchers();
  },

  // Get color scheme from local storage
  get schemeFromLocalStorage() {
    return (
      window.localStorage?.getItem(this.localStorageKey) ??
      this.preferredColorScheme
    );
  },

  // Preferred color scheme
  get preferredColorScheme() {
    return window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light";
  },

  // Init switchers
  initSwitchers() {
    const toggleSwitch = document.getElementById("theme-toggle");

    toggleSwitch.addEventListener(
      "change",
      (e) => {
        // e.preventDefault();
        if (e.target.checked) {
          this.scheme = "dark";
        } else {
          this.scheme = "light";
        }
      },
      false,
    );
  },

  // Set scheme
  set scheme(scheme) {
    this._scheme = scheme;
    this.applyScheme();
    this.schemeToLocalStorage();
  },

  // Get scheme
  get scheme() {
    return this._scheme;
  },

  // Apply scheme
  applyScheme() {
    document
      .querySelector("html")
      ?.setAttribute(this.rootAttribute, this.scheme);
  },

  // Store scheme to local storage
  schemeToLocalStorage() {
    window.localStorage?.setItem(this.localStorageKey, this.scheme);
  },
};

// Init
themeSwitcher.init();
