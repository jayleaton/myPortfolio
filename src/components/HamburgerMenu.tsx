"use client"

import React, { useState } from "react"
import { FiMenu, FiX } from "react-icons/fi"
import { Button } from "@/components/Button"
import { FaGithub, FaLinkedin } from "react-icons/fa"
import { BiLogoUpwork } from "react-icons/bi"

export const HamburgerMenu = () => {
  const [isOpen, setIsOpen] = useState(true)

  function toggleMenu() {
    setIsOpen(!isOpen)
  }

  return (
    <div className="w-full mx-auto z-20">
      {isOpen ? (
        <div className={`absolute right-5 top-5 hover:text-font_highlight text-font_main z-10`}>
          <FiMenu size={35} onClick={toggleMenu} className="hover:cursor-pointer" />
        </div>
      ) : (
        <div className="absolute w-full right-0 top-16 z-10">
          <div className="w-3/4 mx-auto h-auto p-4 pb-8 rounded-xl z-40 bg-background text-font_light opacity-95">
            <FiX size={45} onClick={toggleMenu} className="hover:text-font_highlight hover:cursor-pointer" />
            <div className="w-full font-code gap-2 text-[18px] mx-auto flex flex-col justify-center items-center">
              <Button path={"#Companies"} click={toggleMenu} className="hover:text-font_highlight">
                Companies
              </Button>
              <Button path={"#Projects"} click={toggleMenu} className="hover:text-font_highlight">
                Projects
              </Button>
              <Button path={"#Skills"} click={toggleMenu} className="hover:text-font_highlight">
                Skills
              </Button>
              <Button path={"#About"} click={toggleMenu} className="hover:text-font_highlight">
                About
              </Button>
              <Button path="resume.pdf" target="_blank" click={toggleMenu} className="hover:text-font_highlight">
                Resume
              </Button>
              <Button path="https://github.com/oldbettie" target="_blank">
                <div className="flex justify-center items-center gap-3">
                  <FaGithub className="text-2xl" />
                  Github
                </div>
              </Button>
              <Button path="https://www.linkedin.com/in/jayj-eaton" target="_blank">
                <div className="flex justify-center items-center gap-3">
                  <FaLinkedin className="text-2xl" />
                  Linkedin
                </div>
              </Button>
              <Button path="https://www.upwork.com/freelancers/~01bb2932f2d3e396a6" target="_blank">
                <div className="flex justify-center items-center gap-3">
                  <BiLogoUpwork className="text-3xl" />
                  Upwork
                </div>
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
