package main

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"
)

func AutoDiscover(allFiles []string) ([]string, map[string]string) {
	list := discoverByExplicitProjectType()

	if len(list) == 0 {
		list = discoverByExistingPatterns(allFiles)
	}

	return maps.Keys(list), list
}

func readFromSelections(allFiles []string, opts GIGenOptions) ([]string, map[string]string) {
	var answer bool
	if opts.AutoDiscover.HasValue {
		answer = opts.AutoDiscover.Value
	} else {
		answer = askDiscovery()
	}

	if !answer {
		return readAllFiles(allFiles)
	}

	return AutoDiscover(allFiles)
}

func readAllFiles(allFiles []string) ([]string, map[string]string) {
	baseNames := []string{}
	for _, fn := range allFiles {
		basename := filepath.Base(fn)
		langName := basename[:strings.Index(basename, ".")]
		baseNames = append(baseNames, langName)
	}
	return baseNames, getAllFiles(allFiles)
}

func getAllFiles(allFiles []string) map[string]string {
	files := make(map[string]string)

	for _, filename := range allFiles {
		contents := ReadFile(filename)
		basename := filepath.Base(filename)
		langName := basename[:strings.Index(basename, ".")]

		files[langName] = contents
	}

	return files
}

func discoverByExistingPatterns(allFiles []string) map[string]string {
	files := make(map[string]string)

	for _, filename := range allFiles {
		contents := ReadFile(filename)
		basename := filepath.Base(filename)
		langName := basename[:strings.Index(basename, ".")]

		if findPatternFileMatches(contents) {
			files[langName] = contents
		}
	}
	return files
}

func discoverByExplicitProjectType() map[string]string {
	wd, err := os.Getwd()
	HandleErr(err)

	discoveryMap := make(map[string]string)

	// Common workspace files
	discoveryMap["*.uproject"] = "UnrealEngine"
	discoveryMap["Gemfile"] = "Ruby"
	discoveryMap["Jenkinsfile"] = "JENKINS_HOME"
	discoveryMap["[tj]sconfig.json"] = "Node"
	discoveryMap["_config.ya?ml"] = "Jekyll"
	discoveryMap["app/manifests/AndroidManifest.xml"] = "Android"
	discoveryMap["bin/rails"] = "Rails"
	discoveryMap["composer.json"] = "Composer"
	discoveryMap["config/boot.rb"] = "Rails"
	discoveryMap["go.{mod,sum}"] = "Go"
	discoveryMap["jobs"] = "JENKINS_HOME"
	discoveryMap["package.json"] = "Node"
	discoveryMap["pubspec.ya?ml"] = "Dart"

	// Extensions
	discoveryMap["*.agda"] = "Agda"
	discoveryMap["*.al"] = "AL"
	discoveryMap["*.as"] = "Actionscript"
	discoveryMap["*.dart"] = "Dart"
	discoveryMap["*.dm"] = "DM"
	discoveryMap["*.el"] = "Elisp"
	discoveryMap["*.elm"] = "Elm"
	discoveryMap["*.gcov"] = "Gcov"
	discoveryMap["*.go"] = "Go"
	discoveryMap["*.godot"] = "Godot"
	discoveryMap["*.gradle"] = "Gradle"
	discoveryMap["*.java"] = "Java"
	discoveryMap["*.jl"] = "Julia"
	discoveryMap["*.js"] = "Node"
	discoveryMap["*.lvproj"] = "LabVIEW"
	discoveryMap["*.nim"] = "Nim"
	discoveryMap["*.opa"] = "Opa"
	discoveryMap["*.pde"] = "Processing"
	discoveryMap["*.rb"] = "ChefCookbook"
	discoveryMap["*.sass"] = "Sass"
	discoveryMap["*.swift"] = "Swift"
	discoveryMap["*.unity"] = "Unity"
	discoveryMap["*.zep"] = "Zephir"
	discoveryMap["*.{adb,ada,ads}"] = "Ada"
	discoveryMap["*.{c,cats,h,idc,w}"] = "C"
	discoveryMap["*.{c,h}"] = "Sdcc"
	discoveryMap["*.{cfm,cfc}"] = "CFWheels"
	discoveryMap["*.{clj,boot,cl2,cljc,cljs,cljs.hl,cljscm,cljx,hic}"] = "Clojure"
	discoveryMap["*.{clj,cljs,edn}"] = "Leiningen"
	discoveryMap["*.{cls,trigger,page,component}"] = "ForceDotCom"
	discoveryMap["*.{cmake,cmake.in}"] = "CMake"
	discoveryMap["*.{coq,v}"] = "Coq"
	discoveryMap["*.{cpp,c++,cc,cp,cxx,h,h++,hh,hpp,hxx,inc,inl,ipp,tcc,tpp}"] = "C++"
	discoveryMap["*.{cpp,h}"] = "FlaxEngine"
	discoveryMap["*.{cpp,h}"] = "Qt"
	discoveryMap["*.{cpp,h}"] = "ROS"
	discoveryMap["*.{cpp,h}"] = "TwinCAT3"
	discoveryMap["*.{ctp,php}"] = "CakePHP"
	discoveryMap["*.{cu,cuh}"] = "CUDA"
	discoveryMap["*.{d,di}"] = "D"
	discoveryMap["*.{erl,es,escript,hrl,xrl,yrl}"] = "Erlang"
	discoveryMap["*.{ex,exs}"] = "Elixir"
	discoveryMap["*.{f90,f,f03,f08,f77,f95,for,fpp}"] = "Fortran"
	discoveryMap["*.{fmb,pll,mmb}"] = "OracleForms"
	discoveryMap["*.{fy,fancypack}"] = "Fancy"
	discoveryMap["*.{groovy,java}"] = "Grails"
	discoveryMap["*.{hs,hsc}"] = "Haskell"
	discoveryMap["*.{idr,lidr}"] = "Idris"
	discoveryMap["*.{java,kt,xml}"] = "Android"
	discoveryMap["*.{java,scala}"] = "PlayFramework"
	discoveryMap["*.{java,xml}"] = "GWT"
	discoveryMap["*.{java,xml}"] = "SeamGen"
	discoveryMap["*.{js,ts,xml}"] = "AppceleratorTitanium"
	discoveryMap["*.{js,ts}"] = "ExtJs"
	discoveryMap["*.{json,hcl}"] = "Packer"
	discoveryMap["*.{kt,ktm,kts}"] = "Kotlin"
	discoveryMap["*.{lisp,lsp}"] = "CommonLisp"
	discoveryMap["*.{lua,fcgi,nse,pd_lua,rbxs,wlua}"] = "Lua"
	discoveryMap["*.{ly,ily}"] = "Lilypond"
	discoveryMap["*.{m,h}"] = "Objective-C"
	discoveryMap["*.{m,moo}"] = "Mercury"
	discoveryMap["*.{md,html}"] = "GitBook"
	discoveryMap["*.{ml,eliom,eliomi,ml4,mli,mll,mly}"] = "OCaml"
	discoveryMap["*.{mus,xml}"] = "Finale"
	discoveryMap["*.{pas,dpr,dpk,dfm,dproj}"] = "Delphi"
	discoveryMap["*.{php,html,css,js}"] = "CodeIgniter"
	discoveryMap["*.{php,html,css,js}"] = "Concrete5"
	discoveryMap["*.{php,html,css,js}"] = "CraftCMS"
	discoveryMap["*.{php,html,css,js}"] = "Drupal"
	discoveryMap["*.{php,html,css,js}"] = "ExpressionEngine"
	discoveryMap["*.{php,html,css,js}"] = "FuelPHP"
	discoveryMap["*.{php,html,css,js}"] = "Joomla"
	discoveryMap["*.{php,html,css,js}"] = "Kohana"
	discoveryMap["*.{php,html,css,js}"] = "Laravel"
	discoveryMap["*.{php,html,css,js}"] = "LemonStand"
	discoveryMap["*.{php,html,css,js}"] = "Lithium"
	discoveryMap["*.{php,html,css,js}"] = "Magento"
	discoveryMap["*.{php,html,css,js}"] = "MetaProgrammingSystem"
	discoveryMap["*.{php,html,css,js}"] = "Nanoc"
	discoveryMap["*.{php,html,css,js}"] = "OpenCart"
	discoveryMap["*.{php,html,css,js}"] = "Phalcon"
	discoveryMap["*.{php,html,css,js}"] = "Plone"
	discoveryMap["*.{php,html,css,js}"] = "Prestashop"
	discoveryMap["*.{php,html,css,js}"] = "Qooxdoo"
	discoveryMap["*.{php,html,css,js}"] = "Raku"
	discoveryMap["*.{php,html,css,js}"] = "Scrivener"
	discoveryMap["*.{php,html,css,js}"] = "Stella"
	discoveryMap["*.{php,html,css,js}"] = "SugarCRM"
	discoveryMap["*.{php,html,css,js}"] = "Symfony"
	discoveryMap["*.{php,html,css,js}"] = "SymphonyCMS"
	discoveryMap["*.{php,html,css,js}"] = "Textpattern"
	discoveryMap["*.{php,html,css,js}"] = "Typo3"
	discoveryMap["*.{php,html,css,js}"] = "VVVV"
	discoveryMap["*.{php,html,css,js}"] = "WordPress"
	discoveryMap["*.{php,html,css,js}"] = "Yeoman"
	discoveryMap["*.{php,html,css,js}"] = "Yii"
	discoveryMap["*.{php,html,css,js}"] = "ZendFramework"
	discoveryMap["*.{pkgbuild,install}"] = "ArchLinuxPackages"
	discoveryMap["*.{pl,al,cgi,fcgi,perl,ph,plx,pm,pod,psgi,t}"] = "Perl"
	discoveryMap["*.{pxp,ipf}"] = "IGORPro"
	discoveryMap["*.{py,bzl,cgi,fcgi,gyp,lmi,pyde,pyp,pyt,pyw,rpy,tac,wsgi,xpy}"] = "Python"
	discoveryMap["*.{py,html,css,js}"] = "TurboGears2"
	discoveryMap["*.{r,rd,rsx}"] = "R"
	discoveryMap["*.{rb,builder,fcgi,gemspec,god,irbrc,jbuilder,mspec,pluginspec,podspec,rabl,rake,rbuild,rbw,rbx,ru,ruby,thor,watchr}"] = "Ruby"
	discoveryMap["*.{rb,erb}"] = "RhodesRhomobile"
	discoveryMap["*.{rkt,rktd,rktl,scrbl}"] = "Racket"
	discoveryMap["*.{rs,rs.in}"] = "Rust"
	discoveryMap["*.{scala,sbt,sc}"] = "Scala"
	discoveryMap["*.{sch,brd,kicad_pcb}"] = "KiCad"
	discoveryMap["*.{sch,brd}"] = "Eagle"
	discoveryMap["*.{scm,sld,sls,sps,ss}"] = "Scheme"
	discoveryMap["*.{skp,rb}"] = "SketchUp"
	discoveryMap["*.{sln,csproj}"] = "VisualStudio"
	discoveryMap["*.{st,cs}"] = "Smalltalk"
	discoveryMap["*.{tex,aux,bbx,bib,cbx,cls,dtx,ins,lbx,ltx,mkii,mkiv,mkvi,sty,toc}"] = "TeX"
	discoveryMap["*.{tf,hcl}"] = "Terraform"
	discoveryMap["*.{wscript,py}"] = "Waf"
	discoveryMap["*.{xml,java}"] = "Maven"
	discoveryMap["*.{xojo_code,xojo_menu,xojo_report,xojo_script,xojo_toolbar,xojo_window}"] = "Xojo"
	discoveryMap["*.{yaml,yml}"] = "AppEngine"
	discoveryMap["configure.ac"] = "Autotools"

	results := make(map[string]string)

	for _, key := range maps.Keys(discoveryMap) {
		langName := discoveryMap[key]
		ignoreFile := filepath.Join(GetCacheDir(), langName+".gitignore")
		checkFile := filepath.Join(wd, key)

		_, keyExists := results[langName]
		if !keyExists && GlobExists(checkFile) {
			results[langName] = ReadFile(ignoreFile)
		}
	}

	return results
}
