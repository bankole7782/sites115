# How to create a blog
**Date: 22nd February 2022**


## Solution 1

Use **Markdown**.

Its a very easy to use markdown and convert it to html in your code. A quick search on Google would provide this.

This approach do not use a database (which are not easy to scale). This approach can be very scalable. This approach can be quite cheap
since you don't need to pay for a database.

This is the approach I used for this blog.

### Problems with this approach.
It is not very comfortable for a situation where there exists multiple editors in a blog.


## Solution 2

Use a **Javascript Editor Library** and a **Database**

A quick search on google can get you many editors, some paid some free. I would like to recommend [trix editor](https://trix-editor.org/) .
I have used this on a project.

In this method you would have to do the saving and retrieval of blog posts to a database.

This approach is easier to the user.

### Problems with this approach
1.  A database is not easily autoscaled and thus this approach can generally become slow.

1.  You would need to pay for both web servers and a database server. Thereby being the expensive solution.
