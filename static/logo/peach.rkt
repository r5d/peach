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
(define (moon-drawing size)
  (let* ((moon (new dc-path%))
         (inner-arc-x (/ (* size -50) 500.0))
         (inner-arc-y (* -1 inner-arc-x))
         (inner-arc-size (/ (* size 320) 500.0))
         (outer-arc-x (/ (* size 5) 500.0))
         (outer-arc-y outer-arc-x)
         (outer-arc-size (/ (* size 490) 500.0)))
    (send moon arc inner-arc-x inner-arc-y inner-arc-size inner-arc-size 1.57 4.36 #f)
    (send moon arc outer-arc-x outer-arc-y outer-arc-size outer-arc-size 3.54 2.20 #t)
    moon))


;; Draws the peach logo in a bitmap and returns the bitmap.
(define (draw-logo size)
  (let* ((target (make-bitmap size size))
         (dc (new bitmap-dc% (bitmap target))))
    (setup-dc dc)
    (send dc set-brush "black" 'solid)
    (send dc draw-path (moon-drawing size))
    target))


;; Exports the logo into PNG.
(define (png-export logo size)
  (send logo save-file (format "peach-~s.png" size) 'png))


;; Makes a sizexsize peach logo and exports it to a PNG.
(define (make-logo size)
  (let ((logo-bitmap (draw-logo size)))
    (make-object image-snip% logo-bitmap)  ;Render logo in racket shell.
    (png-export logo-bitmap size)))        ;Export logo to a PNG file.


;; Make logos in different sizes.
(begin
  (let ((sizes '(58 76 80 87 114 120 152 167 120 180)))
    (map make-logo sizes)))
