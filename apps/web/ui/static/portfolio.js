// The portfolio page's only script besides Datastar (CSP: 'self' + nonce;
// no inline styles, so state is classes and data attributes). The page is
// complete without it: the server renders every word and number.
// - decode: lights the hero's .pf-tok spans one by one at the selected
//   run's published rate, fills the readout, mirrors the run wherever the
//   page shows it and sets --pf-rate for the stack figure's packets. The
//   runs are the [data-run] tabs on the local AI page.
// - roll: [data-odo] follows data-target (Datastar keeps it on the live
//   $followers signal) with one rolling column per digit.
// - grow: [data-reveal] gets .pf-in when it scrolls into view (bars).
// - menu: closes the <details data-menu> on a choice, Escape or a click away.
// prefers-reduced-motion: the decode shows its final state, numbers swap.
const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
const root = document.documentElement
const each = (selector, apply, scope = document) => scope.querySelectorAll(selector).forEach(apply)
const setText = (selector, text) => each(selector, (element) => (element.textContent = text))

function decodeHero() {
  const hero = document.querySelector('[data-hero]')
  const tokens = hero ? [...hero.querySelectorAll('.pf-tok')] : []
  const tabs = [...document.querySelectorAll('[data-run]')]
  // Hand the pre-paint ghost over to the hero, where the tokens light up.
  hero?.classList.add('pf-decoding')
  root.classList.remove('pf-decode')
  if (!hero || tokens.length === 0) return

  let rate = Number(tabs[0]?.dataset.rate) || 80
  let frame = 0
  setText('[data-run-total]', String(tokens.length))

  const paint = (count, seconds) => {
    tokens.forEach((token, index) => {
      token.classList.toggle('pf-on', index < count)
      token.classList.toggle('pf-now', index === count - 1 && count < tokens.length)
    })
    setText('[data-run-done]', String(count))
    setText('[data-run-time]', seconds.toFixed(2))
  }

  const play = () => {
    cancelAnimationFrame(frame)
    const total = tokens.length
    if (reduced) return paint(total, total / rate)
    const start = performance.now()
    const step = (now) => {
      const seconds = Math.max(0, (now - start) / 1000)
      const count = Math.min(total, Math.floor(seconds * rate) + 1)
      if (count >= total) return paint(total, total / rate)
      paint(count, seconds)
      frame = requestAnimationFrame(step)
    }
    frame = requestAnimationFrame(step)
  }

  const select = (tab) => {
    rate = Number(tab.dataset.rate) || rate
    for (const other of tabs) other.setAttribute('aria-pressed', String(other === tab))
    setText('[data-run-rate]', tab.dataset.rate)
    setText('[data-run-name]', tab.dataset.name)
    setText('[data-run-label]', tab.dataset.label)
    each('[data-run-source]', (link) => link.setAttribute('href', tab.dataset.source))
    root.style.setProperty('--pf-rate', String(rate))
    play()
  }

  for (const tab of tabs) tab.addEventListener('click', () => select(tab))
  each('[data-replay]', (button) => button.addEventListener('click', play))
  each('[data-show-tokens]', (button) =>
    button.addEventListener('click', () => {
      const shown = button.getAttribute('aria-pressed') !== 'true'
      button.setAttribute('aria-pressed', String(shown))
      hero.classList.toggle('pf-show-tokens', shown)
    }),
  )
  play()
}

// rollTo rebuilds the odometer for target: each digit is a 0-9 strip that
// starts where the previous number left that column and rolls to its digit.
function rollTo(element, target) {
  if (!Number.isFinite(target)) return
  const text = Math.round(target).toLocaleString('en-US')
  if (reduced) {
    element.textContent = text
    return
  }
  const previous = [...element.querySelectorAll('.pf-odo-strip')].map((strip) => strip.dataset.digit)
  const strips = []
  let place = text.replace(/\D/g, '').length // digits still to place, so columns align from the right
  element.textContent = ''
  for (const character of text) {
    if (!/\d/.test(character)) {
      element.append(character)
      continue
    }
    place -= 1
    const digit = document.createElement('span')
    digit.className = 'pf-odo-digit'
    const strip = document.createElement('span')
    strip.className = 'pf-odo-strip'
    strip.textContent = '0\n1\n2\n3\n4\n5\n6\n7\n8\n9'
    strip.dataset.digit = previous[previous.length - 1 - place] ?? '0'
    digit.append(strip)
    element.append(digit)
    strips.push([strip, character])
  }
  // Two frames: the start positions paint first, then the strips roll.
  requestAnimationFrame(() =>
    requestAnimationFrame(() => {
      for (const [strip, character] of strips) strip.dataset.digit = character
    }),
  )
}

function watchOdometers() {
  const sync = (element) => rollTo(element, Number(element.dataset.target))
  const changed = new MutationObserver((records) => {
    for (const record of records) sync(record.target)
  })
  each('[data-odo]', (element) => {
    changed.observe(element, { attributes: true, attributeFilter: ['data-target'] })
    if (element.dataset.target) sync(element)
  })
}

function reveal() {
  const targets = document.querySelectorAll('[data-reveal]')
  if (reduced || !('IntersectionObserver' in window)) {
    for (const element of targets) element.classList.add('pf-in')
    return
  }
  const observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (!entry.isIntersecting) continue
        entry.target.classList.add('pf-in')
        observer.unobserve(entry.target)
      }
    },
    { rootMargin: '0px 0px -10% 0px' },
  )
  for (const element of targets) observer.observe(element)
}

function menu() {
  each('[data-menu]', (details) => {
    const close = () => details.removeAttribute('open')
    each('a', (link) => link.addEventListener('click', close), details)
    document.addEventListener('keydown', (event) => {
      if (event.key !== 'Escape' || !details.open) return
      close()
      details.querySelector('summary')?.focus()
    })
    document.addEventListener('click', (event) => {
      if (details.open && !details.contains(event.target)) close()
    })
  })
}

decodeHero()
watchOdometers()
reveal()
menu()
