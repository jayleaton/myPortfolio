export function Footer() {
  return (
    <footer className="mt-12 py-20">
      <div className="section-panel rounded-[32px] border border-white/10 px-6 py-10 text-center shadow-[0_30px_80px_rgba(2,6,23,0.25)] md:px-8">
        <h3 className="font-console text-2xl tracking-wide text-font_light md:text-[28px]">
          Let&apos;s build something durable and useful together.
        </h3>
        <p className="mt-4 font-interface text-font_main md:text-lg">
          JJ Eaton{" "}
          <a className="underline transition-colors hover:text-font_highlight" href="mailto: jayleaton@gmail.com">
            jayleaton@gmail.com
          </a>
        </p>
      </div>
    </footer>
  )
}
