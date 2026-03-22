"use client"

import React, { useEffect, useState } from "react"
import { Button } from "@/components/Button"
import { UseMediaQuery } from "@/hooks/useMediaQuery"
import { FaGithub, FaLinkedin } from "react-icons/fa"
import { BiLogoUpwork } from "react-icons/bi"
import { HamburgerMenu } from "@/components/HamburgerMenu"

const getIsDesktop = (): boolean => {
  return UseMediaQuery(900)
}

enum NavSections {
  "COMPANIES" = "Companies",
  "PROJECTS" = "Projects",
  "SKILLS" = "Skills",
  "ABOUT" = "About",
  "RESUME" = "Resume",
}

export function Navigation() {
  const isDesktop = getIsDesktop()
  const [activeSection, setActiveSection] = useState<NavSections>(NavSections.COMPANIES)

  useEffect(() => {
    const sectionMap: Record<string, NavSections> = {
      Companies: NavSections.COMPANIES,
      Projects: NavSections.PROJECTS,
      Skills: NavSections.SKILLS,
      About: NavSections.ABOUT,
    }

    const sectionObserver = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            const nextSection = sectionMap[entry.target.id]
            if (nextSection) setActiveSection(nextSection)
          }
        })
      },
      { threshold: 0.5 }
    )

    const sections = document.querySelectorAll("#Companies, #Projects, #Skills, #About")
    sections.forEach((section) => {
      sectionObserver.observe(section)
    })

    return () => {
      sectionObserver.disconnect()
    }
  }, [])

  if (isDesktop) {
    return (
      <div className="sticky top-0 z-30 pt-5">
        <div className="flex w-full items-center justify-between rounded-full border border-white/10 bg-slate-950/35 px-4 py-3 shadow-[0_20px_80px_rgba(2,6,23,0.35)] backdrop-blur-xl xl:px-6">
          <h1 className="font-interface text-2xl font-bold tracking-tight text-white xl:text-3xl">
            JJ <span className="text-font_highlight">Eaton</span>
          </h1>

          <div className="w-fit">
            <div className="flex items-center space-x-3 pl-5 text-[15px] lg:space-x-5">
              <Button
                path="#Companies"
                className={`rounded-full px-3 py-2 font-code transition-all duration-300 hover:bg-white/5 hover:text-font_highlight ${
                  activeSection === NavSections.COMPANIES
                    ? "bg-font_highlight/10 text-font_highlight"
                    : "text-font_light"
                }`}
                click={() => setActiveSection(NavSections.COMPANIES)}
              >
                {activeSection === NavSections.COMPANIES ? `< Companies />` : `Companies`}
              </Button>
              <Button
                path="#Projects"
                className={`rounded-full px-3 py-2 font-code transition-all duration-300 hover:bg-white/5 hover:text-font_highlight ${
                  activeSection === NavSections.PROJECTS ? "bg-font_highlight/10 text-font_highlight" : "text-font_light"
                }`}
                click={() => setActiveSection(NavSections.PROJECTS)}
              >
                {activeSection === NavSections.PROJECTS ? `< Projects />` : `Projects`}
              </Button>
              <Button
                path="#Skills"
                className={`rounded-full px-3 py-2 font-code transition-all duration-300 hover:bg-white/5 hover:text-font_highlight ${
                  activeSection === NavSections.SKILLS ? "bg-font_highlight/10 text-font_highlight" : "text-font_light"
                }`}
                click={() => setActiveSection(NavSections.SKILLS)}
              >
                {activeSection === NavSections.SKILLS ? `< Skills />` : `Skills`}
              </Button>
              <Button
                path="#About"
                className={`rounded-full px-3 py-2 font-code transition-all duration-300 hover:bg-white/5 hover:text-font_highlight ${
                  activeSection === NavSections.ABOUT ? "bg-font_highlight/10 text-font_highlight" : "text-font_light"
                }`}
                click={() => setActiveSection(NavSections.ABOUT)}
              >
                {activeSection === NavSections.ABOUT ? `< About />` : `About`}
              </Button>
              <Button
                path="resume.pdf"
                target="_blank"
                className={`rounded-full px-3 py-2 font-code transition-all duration-300 hover:bg-white/5 hover:text-font_highlight ${
                  activeSection === NavSections.RESUME ? "bg-font_highlight/10 text-font_highlight" : "text-font_light"
                }`}
                click={() => setActiveSection(NavSections.RESUME)}
              >
                {activeSection === NavSections.RESUME ? `< Resume />` : `Resume`}
              </Button>
            </div>
          </div>

          <div className="flex gap-2 xl:gap-3">
            <Button
              path="https://github.com/oldbettie"
              target="_blank"
              className="rounded-full border border-white/10 px-3 py-2 text-font_light transition-all duration-300 hover:-translate-y-0.5 hover:border-white/20 hover:text-white"
            >
              <div className="flex items-center justify-center gap-2">
                <FaGithub className="text-lg" />
                <span className="hidden lg:inline">Github</span>
              </div>
            </Button>
            <Button
              path="https://www.linkedin.com/in/jayj-eaton"
              target="_blank"
              className="rounded-full border border-white/10 px-3 py-2 text-font_light transition-all duration-300 hover:-translate-y-0.5 hover:border-blue-400/20 hover:text-blue-300"
            >
              <div className="flex items-center justify-center gap-2">
                <FaLinkedin className="text-lg" />
                <span className="hidden lg:inline">Linkedin</span>
              </div>
            </Button>
            <Button
              path="https://www.upwork.com/freelancers/~01bb2932f2d3e396a6"
              target="_blank"
              className="rounded-full border border-white/10 px-3 py-2 text-font_light transition-all duration-300 hover:-translate-y-0.5 hover:border-green-400/20 hover:text-green-300"
            >
              <div className="flex items-center justify-center gap-2">
                <BiLogoUpwork className="text-xl" />
                <span className="hidden lg:inline">Upwork</span>
              </div>
            </Button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="sticky top-0 z-30 pt-4">
      <div className="flex w-full justify-between rounded-full border border-white/10 bg-slate-950/40 px-4 py-3 shadow-[0_20px_80px_rgba(2,6,23,0.35)] backdrop-blur-xl">
        <h1 className="w-full font-interface text-2xl text-font_main">JJ Eaton</h1>
        <HamburgerMenu />
      </div>
    </div>
  )
}
