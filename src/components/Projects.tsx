"use client"

import { projects } from "@/utils/projects"
import { useState } from "react"
import { Button } from "@/components/Button"
import { FaExternalLinkSquareAlt, FaGithub } from "react-icons/fa"
import Link from "next/link"

export function Projects() {
  const [counter, setCounter] = useState<number>(4)
  return (
    <div id="Projects" className="py-20 relative">
      <div className="mb-12 flex flex-col items-center text-center md:flex-row md:items-end md:justify-between md:text-left">
        <div className="md:w-1/2">
           {/* Placeholder if needed */}
        </div>
        <h2 className="text-center text-[35px] font-bold font-code leading-tight md:text-right md:text-[50px]">
          Latest <br />
          <span className="text-transparent bg-clip-text bg-gradient-to-r from-font_highlight to-blue-500">
            Projects
          </span>
        </h2>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8 pt-5">
        {projects.map((project, num) => {
          if (num < counter) {
            return (
              <div
                key={project.name}
                className="group relative h-full bg-background_highlight/40 backdrop-blur-sm border border-white/5 rounded-2xl p-6 transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl hover:shadow-font_highlight/10 flex flex-col"
              >
                <div className="mb-6 flex flex-col items-center gap-4 text-center md:flex-row md:items-start md:justify-between md:text-left">
                  <h3 className="text-white text-2xl font-code font-bold group-hover:text-font_highlight transition-colors">
                    {project.name}
                  </h3>
                  <div className="flex gap-2">
                    {project.github && (
                      <Button
                        path={project.github}
                        target="_blank"
                        className="p-2 rounded-full hover:bg-white/10 transition-colors text-font_light hover:text-white"
                      >
                        <FaGithub size={24} />
                      </Button>
                    )}
                    <Button
                      className="p-2 rounded-full hover:bg-white/10 transition-colors text-font_highlight hover:text-font_highlight/80"
                      path={project.website}
                      target="_blank"
                    >
                      <FaExternalLinkSquareAlt size={24} />
                    </Button>
                  </div>
                </div>

                <Link href={project.website} target="_blank" className="block w-full overflow-hidden rounded-xl mb-6">
                  <img
                    src={project.image}
                    alt={project.name}
                    className="w-full h-auto transform transition-transform duration-500 group-hover:scale-105"
                  />
                </Link>

                <div className="flex-grow">
                  <p className="mb-6 text-center text-base leading-relaxed text-font_main md:text-left">{project.about}</p>
                </div>

                <div className="pt-4 border-t border-white/5">
                  <div className="flex flex-wrap gap-6 justify-center text-font_light/70">
                    {project.stack.map((Icon, idx) => (
                      <Icon key={idx} size={32} className="hover:text-font_highlight transition-colors" />
                    ))}
                  </div>
                </div>
              </div>
            )
          }
          return null
        })}
      </div>
      <div className="py-10 flex justify-center">
        {counter === 4 && (
          <Button
            className="font-code text-xl border border-font_highlight/30 text-font_highlight px-8 py-3 rounded-full hover:bg-font_highlight/10 transition-all hover:scale-105"
            click={() => setCounter(10)}
          >
            Show More Projects
          </Button>
        )}
        {counter === 10 && (
          <Button
            className="font-code text-xl border border-font_highlight/30 text-font_highlight px-8 py-3 rounded-full hover:bg-font_highlight/10 transition-all hover:scale-105"
            click={() => setCounter(4)}
          >
            Show Less
          </Button>
        )}
      </div>
    </div>
  )
}
