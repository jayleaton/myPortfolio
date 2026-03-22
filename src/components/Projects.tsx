"use client"

import { projects } from "@/utils/projects"
import { useState } from "react"
import { Button } from "@/components/Button"
import { FaExternalLinkSquareAlt, FaGithub } from "react-icons/fa"
import Link from "next/link"

export function Projects() {
  const [counter, setCounter] = useState<number>(4)
  return (
    <section id="Projects" className="section-shell relative scroll-mt-28 py-20">
      <div className="section-grid absolute inset-0 opacity-60" />
      <div className="section-panel relative overflow-hidden rounded-[32px] border border-white/10 px-6 py-10 md:px-8 xl:px-10">
        <div className="mb-12 flex flex-col items-center text-center md:flex-row md:items-end md:justify-between md:text-left">
          <div className="md:w-1/2">
            <p className="section-kicker">Selected product work</p>
            <p className="mt-4 max-w-xl leading-relaxed text-font_main">
              A mix of SaaS products, client builds, and experiments focused on usable interfaces, clear product value,
              and maintainable implementation.
            </p>
          </div>
        </div>

        <div className="grid grid-cols-1 gap-8 pt-5 md:grid-cols-2">
          {projects.map((project, num) => {
            if (num < counter) {
              return (
                <div
                  key={project.name}
                  className="group flex h-full flex-col rounded-[28px] border border-white/10 bg-slate-950/30 p-6 shadow-[0_30px_80px_rgba(2,6,23,0.25)] backdrop-blur-xl transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl hover:shadow-font_highlight/10"
                >
                  <div className="mb-6 flex flex-col items-center gap-4 text-center md:flex-row md:items-start md:justify-between md:text-left">
                    <h3 className="text-2xl font-bold text-white transition-colors group-hover:text-font_highlight font-code">
                      {project.name}
                    </h3>
                    <div className="flex gap-2">
                      {project.github && (
                        <Button
                          path={project.github}
                          target="_blank"
                          className="rounded-full border border-white/10 p-2 text-font_light transition-colors hover:bg-white/10 hover:text-white"
                        >
                          <FaGithub size={22} />
                        </Button>
                      )}
                      <Button
                        className="rounded-full border border-font_highlight/20 p-2 text-font_highlight transition-colors hover:bg-font_highlight/10 hover:text-font_highlight/80"
                        path={project.website}
                        target="_blank"
                      >
                        <FaExternalLinkSquareAlt size={22} />
                      </Button>
                    </div>
                  </div>

                  <Link href={project.website} target="_blank" className="mb-6 block w-full overflow-hidden rounded-[22px] border border-white/10">
                    <img
                      src={project.image}
                      alt={project.name}
                      className="w-full h-auto transform transition-transform duration-500 group-hover:scale-105"
                    />
                  </Link>

                  <div className="flex-grow">
                    <p className="mb-6 text-center text-base leading-relaxed text-font_main md:text-left">{project.about}</p>
                  </div>

                  <div className="border-t border-white/10 pt-4">
                    <div className="flex flex-wrap justify-center gap-4 text-font_light/70">
                      {project.stack.map((Icon, idx) => (
                        <div key={idx} className="rounded-full border border-white/10 bg-white/[0.03] p-3">
                          <Icon size={26} className="transition-colors hover:text-font_highlight" />
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              )
            }
            return null
          })}
        </div>
        <div className="flex justify-center py-10">
          {counter === 4 && (
            <Button
              className="rounded-full border border-font_highlight/30 px-8 py-3 font-code text-xl text-font_highlight transition-all hover:scale-105 hover:bg-font_highlight/10"
              click={() => setCounter(10)}
            >
              Show More Projects
            </Button>
          )}
          {counter === 10 && (
            <Button
              className="rounded-full border border-font_highlight/30 px-8 py-3 font-code text-xl text-font_highlight transition-all hover:scale-105 hover:bg-font_highlight/10"
              click={() => setCounter(4)}
            >
              Show Less
            </Button>
          )}
        </div>
      </div>
    </section>
  )
}
