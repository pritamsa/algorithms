package arrays;


import java.util.*;

public class Bookings {

    public static void main(String[] args) {
        Integer[] arr = {1, 3, 5};
        Integer[] dep = {2, 6, 8};

        Integer[] arr1 = { 510, 492, 448, 465, 277, 302, 615, 380, 127, 51, 682, 328, 329, 513, 291, 507, 38, 106, 300, 15, 118, 273, 469, 85, 111, 262, 470, 256, 351, 109, 5, 542, 532, 336, 535, 103, 353, 91, 306, 584, 597, 12, 346, 195, 268, 400, 305, 18, 154, 404, 612, 502, 392, 511, 63, 221, 86, 335, 150, 553, 536, 213, 198, 672, 543, 210, 180, 274, 531, 194, 653, 463, 74, 429, 598, 90, 16, 648, 178, 607, 363, 156, 523, 32, 657, 467, 489, 652, 35, 571, 318, 193, 175, 218, 578, 187, 378, 222, 617, 231, 345, 339, 483, 98, 39, 215, 259, 34, 265, 348, 600, 41, 666, 606, 639, 540, 616, 170, 361, 275, 586, 276, 357, 638, 266, 121, 79, 171, 152, 332, 388, 426, 83, 94, 254, 537, 28, 487, 133, 99, 286, 123, 80, 183, 452, 330, 119, 71, 634, 556, 544, 624, 296, 49, 69, 243, 419, 323, 189, 403, 472, 61, 287, 440, 140, 242, 299, 415, 435, 420, 522, 216, 320, 229, 425, 619, 283, 281, 122, 601, 3, 342, 604, 466, 627, 325, 565, 206, 182, 134, 462, 358, 432, 105, 416, 441, 486, 233, 56, 47, 524, 77, 289, 680, 65, 88, 391, 204, 314, 362, 504, 226, 602, 308, 356, 304, 248, 569, 434, 663, 31, 538, 136, 298, 153, 151, 327, 563, 166, 64, 640, 533, 628, 143, 11, 301, 211, 527, 474, 408, 224, 374, 679, 234, 260, 473, 596, 464, 550, 177, 644, 191, 21, 551, 446, 449, 100, 200, 149, 168, 625, 529, 589, 212, 48, 636, 442, 173, 641, 534, 668, 125, 96, 176, 297, 521, 285, 315, 659, 655, 197, 137, 421, 455, 436, 108, 386, 75, 165, 558, 632, 110, 13, 365, 396, 341, 37, 411, 370, 439, 196, 272, 549, 326, 645, 245, 343, 24, 630, 36, 621, 59, 385, 496, 417, 355, 444, 608, 588, 371, 575, 303, 228, 613, 360, 525, 599, 188, 40, 162, 223, 457, 205, 6, 515, 247, 288, 130, 594, 45, 631, 331, 585, 667, 52, 4, 409, 626, 237, 618, 554, 44, 381, 662, 33, 43, 26, 561, 225, 433, 568, 271, 128, 82, 478, 8, 656, 620, 238, 514, 60, 399, 132, 30, 541, 570, 279, 526, 664, 89, 577, 322, 241, 120, 424, 148, 519, 290, 481, 643, 450, 637, 516, 203, 509, 671, 506, 23, 104, 278, 476, 418, 395, 614, 494, 282, 311, 495, 562, 669, 471, 393, 158, 174, 453, 312, 93, 642, 364, 270, 209, 27, 573, 84, 497, 406, 501, 138, 114, 675, 567, 479, 389, 603, 230, 214, 217, 17, 402, 376, 373, 482, 87, 595, 141, 142, 220, 416, 454, 129, 219, 284, 397, 240, 135, 359, 352, 484, 184, 257, 383, 456, 164, 530, 159, 235, 499, 252, 670, 677, 19, 517, 9, 208, 503, 227, 654, 76, 78, 398, 232, 518, 255, 58, 310, 167, 1, 488, 623, 498, 161, 244, 249, 124, 20, 490, 581, 81, 250, 451, 280, 491, 295, 459, 115, 67, 431, 423, 338, 354, 528, 673, 477, 192, 548, 377, 559, 201, 622, 447, 236, 101, 593, 555, 251, 139, 144, 545, 547, 458, 674, 172, 55, 319, 324, 485, 14, 387, 579, 337, 582, 660, 155, 651, 181, 7, 344, 591, 384, 317, 678, 475, 50, 57, 574, 520, 552, 113, 505, 372, 676, 661, 292, 592, 66, 564, 53, 112, 412, 428, 334, 185, 610, 665, 349, 68, 605, 126, 590, 160, 169, 207, 394, 261, 131, 681, 258, 190, 500, 2, 580, 107, 269, 493, 163, 264, 368, 146, 468, 576, 413, 650, 539, 422, 560, 390, 316, 460, 10, 309, 73, 54, 572, 375, 92, 307, 443, 25, 407, 253, 46, 347, 609, 410, 199, 202, 70, 430, 587, 427, 382, 438, 62, 658, 147, 22, 366, 267, 566, 186, 635, 445, 512, 263, 629, 42, 461, 102, 633, 313, 29, 369, 239, 340, 437, 116, 508, 480, 583, 367, 72, 321, 179, 157, 557, 293, 117, 414, 379, 405, 649, 246, 546, 401, 294, 350, 647, 333, 646, 97, 145, 95, 611 };


        repeatedNumber(new ArrayList<>(Arrays.asList(arr1)));
        //hotel(new ArrayList(Arrays.asList(arr)), new ArrayList(Arrays.asList(dep)), 1);

    }
    static class Resv {
        String ls;
        public Resv(String ls) {
            this.ls = ls;
        }
    }
    public static boolean hotel(ArrayList<Integer> arrive, ArrayList<Integer> depart, int K) {

        Resv[][] lst = new Resv[31][31];
        int i = 0;

        while (i < arrive.size()) {
            Integer arDate = arrive.get(i) != null ? arrive.get(i) : 0;
            Integer deDate = depart.get(i) != null ? depart.get(i) : 0;

            Resv[] lsA = lst[arDate - 1];
            Resv arrivals = null;
            if (lsA == null) {
                lsA = new Resv[2];
                arrivals = new Resv("");

            } else {
                arrivals = lsA[0];

                if (arrivals == null) {
                    arrivals = new Resv("");
                }

            }
            arrivals.ls += "a";
            lsA[0] = arrivals;
            lst[arDate -1 ] = lsA;


            Resv[] lsD = lst[deDate - 1];
            Resv departs = null;
            if (lsD == null) {
                lsD = new Resv[2];
                departs = new Resv("");

            } else {
                departs = lsD[0];

                if (departs == null) {
                    departs = new Resv("");
                }

            }
            departs.ls += "d";
            lsD[1] = departs;
            lst[deDate-1] = lsD;
            i++;



        }
        int avRooms = K;
        for (int p = 0; p < lst.length; p++) {
            Resv[] reservations =  lst[p];
            if (reservations != null && reservations.length > 0) {
                Resv arrivs = reservations[0];
                Resv depts = reservations[1];
                if (depts != null) {
                    avRooms += depts.ls.trim().length();
                }

                if (arrivs != null) {
                    avRooms -= arrivs.ls.trim().length();
                }
                if (avRooms < 0 || avRooms > K) {
                    return false;
                }
            }


        }
        return avRooms == K;



    }

    public static boolean soln(ArrayList<Integer> arrive, ArrayList<Integer> depart, int K) {
        Collections.sort(arrive);
        Collections.sort(depart);

        int i = 0;
        int j = 0;

        while(i<arrive.size() && j<depart.size()){
            if(arrive.get(i)<depart.get(j)){
                i++;
                K--;
            }else if(arrive.get(i)==depart.get(j)){
                i++;
                j++;
            }else{
                j++;
                K++;
            }

            if(K<0)
                return false;
        }
        return true;
    }

    public static int repeatedNumber(final List<Integer> A) {

        List<Integer> B = new ArrayList<>();
        B.addAll(A);
        int j = 0;
        while(j < B.size()) {

            int idx = (B.get(j)%B.size()) - 1;
            if (idx ==653) {
                System.out.println("");
            }
            B.set(idx, (B.get(idx) + B.size()));
            j++;

        }
        for (int i = 0; i < B.size(); i++) {
            if (B.get(i) / B.size() > 1) {
                return i+1;
            }
        }
        return -1;
    }
}
