;;;; Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
;;;; SPDX-License-Identifier: ISC

#lang racket

(require racket/draw)
(require racket/snip)

;; Sets up settings for the drawing context.
(define (setup-dc dc)
  (send dc set-smoothing 'aligned)
  (send dc set-pen "black" 1 'transparent)
  (send dc set-brush "black" 'solid))

;; Draws a moon.
(define (moon-drawing)
  (let ((moon (new dc-path%)))
    (send moon arc -50 50 320 320 1.57 4.36 #f)
    (send moon arc 1 4 490 490 3.54 2.20 #t)
    moon))

;; Draws the peach logo in a bitmap and returns the bitmap.
(define (draw-logo)
  (let* ((target (make-bitmap 500 500))
         (dc (new bitmap-dc% (bitmap target))))
    (setup-dc dc)
    (send dc draw-path (moon-drawing))
    target))

;; Exports the logo into PNG.
(define (png-export logo)
  (send logo save-file "peach.png" 'png))

;; Peach logo as a bitmap.
(define peach-logo (draw-logo))

;; Render logo in racket shell.
(make-object image-snip% peach-logo)

;; Export logo to a PNG file.
(png-export peach-logo)
