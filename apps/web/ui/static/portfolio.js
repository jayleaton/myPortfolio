// The portfolio page's only script besides Datastar (CSP: 'self' + nonce).
// - [data-count]: animates its text up to data-target (set by Datastar from
//   the live $followers signal) the first time it is visible and on every
//   change after. The server renders the final number, so nothing breaks
//   without JS.
// - [data-telemetry]: decorative readouts that jitter around data-base.
// - [data-reveal]: fades sections in as they scroll into view.
// prefers-reduced-motion: numbers jump, nothing jitters or fades.
const reduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches
const root = document.documentElement
const format = (value) => Math.round(value).toLocaleString('en-US')
const duration = () =>
  Number.parseInt(getComputedStyle(root).getPropertyValue('--pf-count-duration'), 10) || 1200

function countTo(element, target) {
  if (!Number.isFinite(target)) return
  const from = Number(element.dataset.shown ?? 0)
  element.dataset.shown = String(target)
  if (reduced || from === target) {
    element.textContent = format(target)
    return
  }
  const start = performance.now()
  const span = duration()
  const step = (now) => {
    const progress = Math.min(1, (now - start) / span)
    const eased = 1 - (1 - progress) ** 3
    element.textContent = format(from + (target - from) * eased)
    if (progress < 1) requestAnimationFrame(step)
  }
  requestAnimationFrame(step)
}

function watchCounts() {
  const visible = new WeakSet()
  const sync = (element) => {
    if (visible.has(element)) countTo(element, Number(element.dataset.target))
  }
  const seen = new IntersectionObserver((entries) => {
    for (const entry of entries) {
      if (!entry.isIntersecting) continue
      visible.add(entry.target)
      seen.unobserve(entry.target)
      sync(entry.target)
    }
  })
  const changed = new MutationObserver((records) => {
    for (const record of records) sync(record.target)
  })
  for (const element of document.querySelectorAll('[data-count]')) {
    seen.observe(element)
    changed.observe(element, { attributes: true, attributeFilter: ['data-target'] })
  }
}

function jitterTelemetry() {
  if (reduced) return
  const readouts = [...document.querySelectorAll('[data-telemetry]')]
  setInterval(() => {
    for (const element of readouts) {
      const base = Number(element.dataset.base)
      const jitter = Number(element.dataset.jitter)
      const value = base + (Math.random() * 2 - 1) * jitter
      element.textContent = value.toFixed(Number(element.dataset.decimals))
    }
  }, 1400)
}

function reveal() {
  if (reduced || !('IntersectionObserver' in window)) return
  root.classList.add('pf-js')
  const observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (!entry.isIntersecting) continue
        entry.target.classList.add('pf-in')
        observer.unobserve(entry.target)
      }
    },
    { rootMargin: '0px 0px -8% 0px' },
  )
  for (const element of document.querySelectorAll('[data-reveal]')) observer.observe(element)
}

watchCounts()
jitterTelemetry()
reveal()
