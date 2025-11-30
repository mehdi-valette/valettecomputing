class VsCalendarDay extends HTMLElement {
  /** @type {HTMLElement} main container of this custom element */
  #container;

  /** @type {number} */
  #totalMinutes;

  constructor() {
    super();
    this.#container = this.#createContainer();
    this.#totalMinutes = 1440;
  }

  connectedCallback() {
    this.append(this.#container);

    // Set up observer
    this.observer = new MutationObserver(this.onMutation);

    // Watch the Light DOM for child node changes
    this.observer.observe(this, {
      childList: true,
    });
  }

  disconnectedCallback() {
    this.observer?.disconnect();
  }

  /** @type {MutationCallback} */
  onMutation = (mutations) => {
    mutations.forEach((mutation) => {
      this.#container.append(...mutation.addedNodes);

      mutation.addedNodes.forEach((node) => {
        if (!(node instanceof VsPeriod)) return;

        node.setParent(this);
      });
    });
  };

  #createContainer() {
    const container = document.createElement("div");
    container.style.width = "500px";
    container.style.height = "1200px";
    container.style.position = "relative";
    container.style.backgroundColor = "lightgray";

    return container;
  }

  // get the number of pixels per minute
  get pixelStep() {
    return this.#container.scrollHeight / this.#totalMinutes;
  }

  get totalMinutes() {
    return this.#totalMinutes;
  }

  get height() {
    return this.#container.clientHeight;
  }

  get top() {
    return (
      this.#container.getBoundingClientRect().top +
      document.documentElement.scrollTop
    );
  }
}

customElements.define("vs-calendar-day", VsCalendarDay);
