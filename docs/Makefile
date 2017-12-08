.PHONY: clean docs

clean:
	shopt -s nullglob; \
	rm *.aux *.log *.gz *.pdf *.ps *.dvi *.out *.fls *.fdb_latexmk \
		img-src/*.aux img-src/*.log img-src/*.gz img-src/*.pdf img-src/*.ps \
		img-src/*.dvi img-src/*.out img-src/*.fls img-src/*.fdb_latexmk

relay-states.pdf: relay-states.tex
	pdflatex -halt-on-error relay-states.tex; \
	pdflatex -halt-on-error relay-states.tex

docs: relay-states.pdf