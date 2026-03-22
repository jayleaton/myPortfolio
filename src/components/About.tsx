export function About() {
  return (
    <section id="About" className="section-shell relative scroll-mt-28 py-20">
      <div className="section-grid absolute inset-0 opacity-60" />
      <div className="section-panel relative overflow-hidden rounded-[32px] border border-white/10 px-6 py-10 md:px-8 xl:px-10">
        <div className="grid items-center gap-12 lg:grid-cols-2">
          <div className="flex flex-col gap-6 font-interface text-lg leading-relaxed text-font_light">
            <div>
              <p className="section-kicker">Who I am</p>
            </div>
            <p>
              Thanks for checking out the portfolio. I keep it current, but if you are interested in working together
              the fastest route is still{" "}
              <a className="font-bold text-font_highlight hover:underline" href="mailto: jayleaton@gmail.com">
                email
              </a>{" "}
              or one of the platforms linked above. I am available for freelance work and product-focused engineering roles.
            </p>
            <p>
              My specialty is early-stage SaaS architecture. I care about building things the right way the first time,
              with enough structure to scale without forcing a rewrite every few months.
            </p>
            <p>
              If I am not writing code, I am usually working on side projects, in the gym, camping, or mountain biking.
            </p>
            <p className="font-bold text-white">
              I am always open to a quick 10 to 15 minute conversation about how I can help move your product forward.
            </p>
          </div>
          <div className="relative mx-auto group">
            <div className="absolute -inset-2 rounded-[30px] bg-[linear-gradient(135deg,rgba(45,212,191,0.25),rgba(96,165,250,0.2))] blur-2xl opacity-60 transition duration-700 group-hover:opacity-90" />
            <div className="relative overflow-hidden rounded-[28px] border border-white/10 bg-slate-950/40 p-3 backdrop-blur-xl">
              <img
                src="/me.jpg"
                alt="JJ Eaton portrait"
                className="relative max-h-[520px] rounded-[22px] shadow-2xl"
              />
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
